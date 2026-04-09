#!/usr/bin/env bash
set -euo pipefail

API_URL="${API_URL:-http://127.0.0.1:8090}"
EMAIL="${EMAIL:-test@test.com}"
PASSWORD="${PASSWORD:-fsca@test}"
SEED_SUFFIX="${SEED_SUFFIX:-$(date +%s)}"
DECKS_FILE="${DECKS_FILE:-scripts/seed-flashcards.json}"
SEED_USAGE_DAYS="${SEED_USAGE_DAYS:-365}"
SEED_SNAPSHOTS="${SEED_SNAPSHOTS:-12}"
SEED_SESSIONS_PER_DECK="${SEED_SESSIONS_PER_DECK:-24}"

require_bin() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "Missing required binary: $1" >&2
		exit 1
	fi
}

require_bin curl
require_bin jq

if [[ ! -f "$DECKS_FILE" ]]; then
	echo "Deck seed file not found: $DECKS_FILE" >&2
	exit 1
fi

login_payload=$(jq -n --arg identity "$EMAIL" --arg password "$PASSWORD" '{identity: $identity, password: $password}')
login_response=$(curl -sS \
	-H "Content-Type: application/json" \
	-X POST "$API_URL/api/collections/users/auth-with-password" \
	-d "$login_payload")

auth_token=$(echo "$login_response" | jq -r '.token // empty')
user_id=$(echo "$login_response" | jq -r '.record.id // empty')

if [[ -z "$auth_token" || -z "$user_id" ]]; then
	echo "Failed to authenticate user $EMAIL" >&2
	echo "$login_response" >&2
	exit 1
fi

auth_header="Authorization: Bearer $auth_token"

SEED_JSON=$(cat "$DECKS_FILE")
CATEGORIES_JSON=$(echo "$SEED_JSON" | jq -c '.categories // []')
DECKS_JSON=$(echo "$SEED_JSON" | jq -c '.decks // .')
NOW_TS=$(date -u +%s)

deck_card_ids=()
deck_card_correct=()
deck_card_wrong=()
deck_card_answer_type=()
deck_card_levels=()
deck_card_difficulty=()
deck_card_wrong_options=()
deck_card_last_review_ts=()
deck_card_next_review_ts=()
deck_card_review_count=()

if (( SEED_USAGE_DAYS < 1 )); then
	echo "SEED_USAGE_DAYS must be at least 1" >&2
	exit 1
fi

if (( SEED_SNAPSHOTS < 1 )); then
	echo "SEED_SNAPSHOTS must be at least 1" >&2
	exit 1
fi

if (( SEED_SESSIONS_PER_DECK < 1 )); then
	echo "SEED_SESSIONS_PER_DECK must be at least 1" >&2
	exit 1
fi

category_names=()
category_ids=()

get_category_id_by_name() {
	local name="$1"
	for i in "${!category_names[@]}"; do
		if [[ "${category_names[$i]}" == "$name" ]]; then
			echo "${category_ids[$i]}"
			return
		fi
	done
	echo ""
}

create_category() {
	local category_json="$1"
	local name description color payload response category_id

	name=$(echo "$category_json" | jq -r '.name')
	description=$(echo "$category_json" | jq -r '.description // empty')
	color=$(echo "$category_json" | jq -r '.color // empty')

	payload=$(jq -n \
		--arg name "$name" \
		--arg description "$description" \
		--arg color "$color" \
		--arg user_id "$user_id" \
		'{
			name: $name,
			user_id: $user_id,
		} + (if ($description | length) > 0 then {description: $description} else {} end) + (if ($color | length) > 0 then {color: $color} else {} end)')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_categories/records" \
		-d "$payload")

	category_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$category_id" ]]; then
		echo "Failed to create category: $name" >&2
		echo "$response" >&2
		exit 1
	fi

	echo "$category_id"
}

link_deck_to_category() {
	local deck_id="$1"
	local category_id="$2"
	local payload response link_id

	payload=$(jq -n --arg deck_id "$deck_id" --arg category_id "$category_id" '{deck_id: $deck_id, category_id: $category_id}')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_deck_categories/records" \
		-d "$payload")

	link_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$link_id" ]]; then
		echo "Failed to link deck $deck_id to category $category_id" >&2
		echo "$response" >&2
		exit 1
	fi
}

create_deck() {
	local deck_json="$1"
	local title slug_base slug color description payload response deck_id

	title=$(echo "$deck_json" | jq -r '.title')
	slug_base=$(echo "$deck_json" | jq -r '.slug')
	slug="${slug_base}-${SEED_SUFFIX}"
	color=$(echo "$deck_json" | jq -r '.color')
	description=$(echo "$deck_json" | jq -r '.description // empty')

	payload=$(jq -n \
		--arg title "$title" \
		--arg slug "$slug" \
		--arg color "$color" \
		--arg user_id "$user_id" \
		--arg description "$description" \
		'{
			title: $title,
			slug: $slug,
			color: $color,
			user_id: $user_id,
			is_public: false
		} + (if ($description | length) > 0 then {description: $description} else {} end)')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_decks/records" \
		-d "$payload")

	deck_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$deck_id" ]]; then
		echo "Failed to create deck: $title" >&2
		echo "$response" >&2
		exit 1
	fi

	echo "$deck_id|$title|$slug"
}

create_card() {
	local card_json="$1"
	local payload response card_id

	payload=$(echo "$card_json" | jq -c --arg creator_id "$user_id" '{
		question,
		answer_type,
		difficulty_level,
		creator_id: $creator_id,
		answers: (
			if (.answers | type) == "array" and (.answers | length) > 0 then .answers
			elif (.correct_answer | type) == "string" and (.correct_answer | length) > 0 then [ .correct_answer ]
			else []
			end
		),
		correct_answer: (
			if (.correct_answer | type) == "string" and (.correct_answer | length) > 0 then .correct_answer
			elif (.answers | type) == "array" and (.answers | length) > 0 then .answers[0]
			else ""
			end
		)
	}')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_cards/records" \
		-d "$payload")

	card_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$card_id" ]]; then
		echo "Failed to create card" >&2
		echo "$response" >&2
		exit 1
	fi

	echo "$card_id"
}

link_card_to_deck() {
	local card_id="$1"
	local deck_id="$2"
	local payload response link_id

	payload=$(jq -n --arg card_id "$card_id" --arg deck_id "$deck_id" '{card_id: $card_id, deck_id: $deck_id}')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_card_decks/records" \
		-d "$payload")

	link_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$link_id" ]]; then
		echo "Failed to link card $card_id to deck $deck_id" >&2
		echo "$response" >&2
		exit 1
	fi
}

format_epoch() {
	local epoch="$1"

	if date -u -r 0 >/dev/null 2>&1; then
		date -u -r "$epoch" +"%Y-%m-%dT%H:%M:%SZ"
	else
		date -u -d "@$epoch" +"%Y-%m-%dT%H:%M:%SZ"
	fi
}

format_date_only() {
	local epoch="$1"

	if date -u -r 0 >/dev/null 2>&1; then
		date -u -r "$epoch" +"%Y-%m-%d"
	else
		date -u -d "@$epoch" +"%Y-%m-%d"
	fi
}

min_int() {
	if (( $1 < $2 )); then
		echo "$1"
	else
		echo "$2"
	fi
}

max_int() {
	if (( $1 > $2 )); then
		echo "$1"
	else
		echo "$2"
	fi
}

clamp_int() {
	local value="$1"
	local min="$2"
	local max="$3"

	if (( value < min )); then
		value="$min"
	elif (( value > max )); then
		value="$max"
	fi

	echo "$value"
}

random_between() {
	local min="$1"
	local max="$2"

	if (( min >= max )); then
		echo "$min"
		return
	fi

	echo $((min + (RANDOM % (max - min + 1))))
}

reviews_days=()
reviews_counts=()

increment_review_day() {
	local day="$1"

	for i in "${!reviews_days[@]}"; do
		if [[ "${reviews_days[$i]}" == "$day" ]]; then
			reviews_counts[$i]=$((reviews_counts[$i] + 1))
			return
		fi
	done

	reviews_days+=("$day")
	reviews_counts+=(1)
}

get_review_day_count() {
	local day="$1"

	for i in "${!reviews_days[@]}"; do
		if [[ "${reviews_days[$i]}" == "$day" ]]; then
			echo "${reviews_counts[$i]}"
			return
		fi
	done

	echo 0
}

leitner_interval_days() {
	case "$1" in
		1) echo 1 ;;
		2) echo 2 ;;
		3) echo 4 ;;
		4) echo 7 ;;
		5) echo 14 ;;
		6) echo 30 ;;
		7) echo 60 ;;
		*) echo 7 ;;
	esac
}

pick_wrong_answer() {
	local options_json="$1"
	local fallback="${2:-nao sei}"
	local count index

	count=$(jq -r 'length' <<<"$options_json" 2>/dev/null || echo 0)
	if (( count < 1 )); then
		echo "$fallback"
		return
	fi

	index=$((RANDOM % count))
	jq -r --argjson index "$index" '.[$index]' <<<"$options_json"
}

pick_short_answer_wrong() {
	local variants=("nao sei" "nao lembro" "vou revisar" "em branco")
	local index=$((RANDOM % ${#variants[@]}))

	echo "${variants[$index]}"
}

calc_correct_chance() {
	local level="$1"
	local difficulty="$2"
	local answer_type="$3"
	local base penalty chance

	difficulty=$(clamp_int "$difficulty" 1 5)
	base=$((55 + (level * 5)))
	penalty=$(((difficulty - 1) * 6))
	chance=$((base - penalty))

	if [[ "$answer_type" == "short_answer" ]]; then
		chance=$((chance - 10))
	fi

	clamp_int "$chance" 20 95
}

create_leitner_box() {
	local deck_id="$1"
	local card_id="$2"
	local level="$3"
	local last_reviewed_at="$4"
	local next_review_at="$5"
	local payload response box_id

	payload=$(jq -n \
		--arg user_id "$user_id" \
		--arg deck_id "$deck_id" \
		--arg card_id "$card_id" \
		--argjson level "$level" \
		--arg last_reviewed_at "$last_reviewed_at" \
		--arg next_review_at "$next_review_at" \
		'{
			user_id: $user_id,
			deck_id: $deck_id,
			card_id: $card_id,
			level: $level,
			last_reviewed_at: $last_reviewed_at,
			next_review_at: $next_review_at
		}')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_leitner_boxes/records" \
		-d "$payload")

	box_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$box_id" ]]; then
		echo "Failed to create Leitner box for card $card_id (deck $deck_id)" >&2
		echo "$response" >&2
		exit 1
	fi
}

create_progress_snapshot() {
	local deck_id="$1"
	local snapshot_date="$2"
	local total_cards="$3"
	local new_cards="$4"
	local learning_cards="$5"
	local review_cards="$6"
	local mastered_cards="$7"
	local due_cards="$8"
	local average_level="$9"
	local completion_rate="${10}"
	local reviewed_today="${11}"
	local payload response snapshot_id

	payload=$(jq -n \
		--arg user_id "$user_id" \
		--arg deck_id "$deck_id" \
		--arg snapshot_date "$snapshot_date" \
		--argjson total_cards "$total_cards" \
		--argjson new_cards "$new_cards" \
		--argjson learning_cards "$learning_cards" \
		--argjson review_cards "$review_cards" \
		--argjson mastered_cards "$mastered_cards" \
		--argjson due_cards "$due_cards" \
		--argjson average_level "$average_level" \
		--argjson completion_rate "$completion_rate" \
		--argjson reviewed_today "$reviewed_today" \
		'{
			user_id: $user_id,
			deck_id: $deck_id,
			snapshot_date: $snapshot_date,
			total_cards: $total_cards,
			new_cards: $new_cards,
			learning_cards: $learning_cards,
			review_cards: $review_cards,
			mastered_cards: $mastered_cards,
			due_cards: $due_cards,
			average_level: $average_level,
			completion_rate: $completion_rate,
			reviewed_today: $reviewed_today
		}')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_progress_history/records" \
		-d "$payload")

	snapshot_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$snapshot_id" ]]; then
		echo "Failed to create progress snapshot for deck $deck_id" >&2
		echo "$response" >&2
		exit 1
	fi
}

create_snapshot_from_state() {
	local deck_id="$1"
	local snapshot_ts="$2"
	local total_cards="${#deck_card_ids[@]}"
	local new_cards=0
	local learning_cards=0
	local review_cards=0
	local mastered_cards=0
	local due_cards=0
	local level_sum=0
	local average_level=0
	local completion_rate=0
	local snapshot_date
	local reviewed_today
	local snapshot_day

	if (( total_cards == 0 )); then
		return
	fi

	for i in "${!deck_card_ids[@]}"; do
		local level="${deck_card_levels[$i]:-1}"
		local review_count="${deck_card_review_count[$i]:-0}"
		local next_review_ts="${deck_card_next_review_ts[$i]:-0}"

		level_sum=$((level_sum + level))

		if (( review_count == 0 )); then
			new_cards=$((new_cards + 1))
		elif (( level >= 6 )); then
			mastered_cards=$((mastered_cards + 1))
		elif (( level <= 2 )); then
			learning_cards=$((learning_cards + 1))
		else
			review_cards=$((review_cards + 1))
		fi

		if (( review_count > 0 && next_review_ts <= snapshot_ts )); then
			due_cards=$((due_cards + 1))
		fi
	done

	average_level=$((level_sum / total_cards))
	completion_rate=$((mastered_cards * 100 / total_cards))
	snapshot_date=$(format_epoch "$snapshot_ts")
	snapshot_day=$(format_date_only "$snapshot_ts")
	reviewed_today=$(get_review_day_count "$snapshot_day")

	create_progress_snapshot \
		"$deck_id" \
		"$snapshot_date" \
		"$total_cards" \
		"$new_cards" \
		"$learning_cards" \
		"$review_cards" \
		"$mastered_cards" \
		"$due_cards" \
		"$average_level" \
		"$completion_rate" \
		"$reviewed_today"
}

create_study_session() {
	local deck_id="$1"
	local session_id="$2"
	local cards_json="$3"
	local progress="$4"
	local state="$5"
	local payload response session_record_id

	payload=$(jq -n \
		--arg user_id "$user_id" \
		--arg deck_id "$deck_id" \
		--arg session_id "$session_id" \
		--arg state "$state" \
		--argjson progress "$progress" \
		--argjson cards "$cards_json" \
		'{
			user_id: $user_id,
			deck_id: $deck_id,
			session_id: $session_id,
			state: $state,
			progress: $progress,
			cards: $cards
		}')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_study_sessions/records" \
		-d "$payload")

	session_record_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$session_record_id" ]]; then
		echo "Failed to create study session $session_id" >&2
		echo "$response" >&2
		exit 1
	fi
}

create_answer_event() {
	local session_id="$1"
	local card_id="$2"
	local result="$3"
	local user_answer="$4"
	local self_assessed="$5"
	local answered_at="$6"
	local payload response event_id

	payload=$(jq -n \
		--arg session_id "$session_id" \
		--arg card_id "$card_id" \
		--arg user_id "$user_id" \
		--arg result "$result" \
		--arg user_answer "$user_answer" \
		--argjson self_assessed "$self_assessed" \
		--arg answered_at "$answered_at" \
		'{
			session_id: $session_id,
			card_id: $card_id,
			user_id: $user_id,
			result: $result,
			user_answer: $user_answer,
			self_assessed: $self_assessed,
			answered_at: $answered_at
		}')

	response=$(curl -sS \
		-H "Content-Type: application/json" \
		-H "$auth_header" \
		-X POST "$API_URL/api/collections/flashcards_answer_events/records" \
		-d "$payload")

	event_id=$(echo "$response" | jq -r '.id // empty')
	if [[ -z "$event_id" ]]; then
		echo "Failed to create answer event for card $card_id" >&2
		echo "$response" >&2
		exit 1
	fi
}

seed_deck_usage_generated() {
	local deck_id="$1"
	local deck_title="$2"
	local total_cards="${#deck_card_ids[@]}"
	local usage_start_ts deck_start_ts usage_span_days start_gap_days
	local snapshot_interval_days session_interval_days
	local snapshot_ts_list=()
	local session_ts_list=()
	local snapshot_index=0
	local prev_snapshot_ts=0
	local prev_session_ts=0

	if (( total_cards == 0 )); then
		return
	fi

	echo "Seeding usage for deck: $deck_title"


	usage_start_ts=$((NOW_TS - (SEED_USAGE_DAYS * 86400)))
	start_gap_days=$((RANDOM % ((SEED_USAGE_DAYS / 4) + 1)))
	deck_start_ts=$((usage_start_ts + (start_gap_days * 86400)))
	usage_span_days=$((SEED_USAGE_DAYS - start_gap_days))
	if (( usage_span_days < 1 )); then
		usage_span_days=1
	fi

	reviews_days=()
	reviews_counts=()

	for i in "${!deck_card_ids[@]}"; do
		deck_card_levels[$i]=1
		deck_card_last_review_ts[$i]=0
		deck_card_next_review_ts[$i]="$deck_start_ts"
		deck_card_review_count[$i]=0
	done

	if (( SEED_SNAPSHOTS <= 1 )); then
		snapshot_interval_days="$usage_span_days"
	else
		snapshot_interval_days=$((usage_span_days / (SEED_SNAPSHOTS - 1)))
		if (( snapshot_interval_days < 1 )); then
			snapshot_interval_days=1
		fi
	fi

	for ((i = 0; i < SEED_SNAPSHOTS; i++)); do
		local snapshot_day snapshot_ts

		snapshot_day=$((i * snapshot_interval_days))
		if (( snapshot_day > usage_span_days )); then
			snapshot_day="$usage_span_days"
		fi
		snapshot_ts=$((deck_start_ts + (snapshot_day * 86400) + (RANDOM % 21600)))
		if (( snapshot_ts < prev_snapshot_ts )); then
			snapshot_ts=$((prev_snapshot_ts + 3600))
		fi
		if (( snapshot_ts > NOW_TS )); then
			snapshot_ts="$NOW_TS"
		fi
		snapshot_ts_list+=("$snapshot_ts")
		prev_snapshot_ts="$snapshot_ts"
	done

	session_interval_days=$((usage_span_days / SEED_SESSIONS_PER_DECK))
	if (( session_interval_days < 1 )); then
		session_interval_days=1
	fi

	local session_day=0
	for ((i = 0; i < SEED_SESSIONS_PER_DECK; i++)); do
		local session_ts session_time_offset step

		if (( i == 0 )); then
			session_day=$((RANDOM % (session_interval_days + 1)))
		else
			step=$((session_interval_days + (RANDOM % 3)))
			session_day=$((session_day + step))
		fi

		if (( session_day > usage_span_days )); then
			session_day="$usage_span_days"
		fi

		session_time_offset=$((28800 + (RANDOM % 50400)))
		session_ts=$((deck_start_ts + (session_day * 86400) + session_time_offset))
		if (( session_ts > NOW_TS )); then
			session_ts="$NOW_TS"
		fi
		if (( i > 0 && session_ts <= prev_session_ts )); then
			session_ts=$((prev_session_ts + 3600))
		fi
		if (( session_ts > NOW_TS )); then
			session_ts="$NOW_TS"
		fi

		session_ts_list+=("$session_ts")
		prev_session_ts="$session_ts"
	done

	snapshot_index=0
	for i in "${!session_ts_list[@]}"; do
		local session_ts="${session_ts_list[$i]}"
		local session_id="seed-${deck_id}-${i}-${SEED_SUFFIX}"

		while (( snapshot_index < ${#snapshot_ts_list[@]} )); do
			local snapshot_ts="${snapshot_ts_list[$snapshot_index]}"

			if (( snapshot_ts > session_ts )); then
				break
			fi

			create_snapshot_from_state "$deck_id" "$snapshot_ts"
			snapshot_index=$((snapshot_index + 1))
		done

		local due_indices=()
		local new_indices=()
		local review_indices=()

		for card_array_index in "${!deck_card_ids[@]}"; do
			local review_count="${deck_card_review_count[$card_array_index]:-0}"
			local next_review_ts="${deck_card_next_review_ts[$card_array_index]:-0}"

			if (( review_count == 0 )); then
				new_indices+=("$card_array_index")
			elif (( next_review_ts <= session_ts )); then
				due_indices+=("$card_array_index")
			else
				review_indices+=("$card_array_index")
			fi
		done

		local session_min session_max session_size
		session_min=$(min_int 4 "$total_cards")
		session_max=$(min_int 12 "$total_cards")
		if (( session_min < 2 )); then
			session_min="$total_cards"
		fi
		session_size=$(random_between "$session_min" "$session_max")
		if (( session_size < 1 )); then
			session_size=1
		fi

		local session_indices=()
		local target_due target_new
		local due_count="${#due_indices[@]}"
		local new_count="${#new_indices[@]}"

		target_due=$((session_size * 60 / 100))
		if (( target_due < 1 && due_count > 0 )); then
			target_due=1
		fi
		target_due=$(min_int "$target_due" "$due_count")

		for ((j = 0; j < target_due; j++)); do
			local idx=$((RANDOM % ${#due_indices[@]}))
			session_indices+=("${due_indices[$idx]}")
			unset "due_indices[$idx]"
			due_indices=(${due_indices[@]-})
		done

		target_new=$((session_size / 4))
		if (( target_new < 1 && new_count > 0 )); then
			target_new=1
		fi
		target_new=$(min_int "$target_new" "$new_count")

		for ((j = 0; j < target_new; j++)); do
			local idx=$((RANDOM % ${#new_indices[@]}))
			session_indices+=("${new_indices[$idx]}")
			unset "new_indices[$idx]"
			new_indices=(${new_indices[@]-})
		done

		local remaining=$((session_size - ${#session_indices[@]}))
		if (( remaining > 0 )); then
			local candidate_indices=(${due_indices[@]-} ${new_indices[@]-} ${review_indices[@]-})
			for ((j = 0; j < remaining && ${#candidate_indices[@]} > 0; j++)); do
				local idx=$((RANDOM % ${#candidate_indices[@]}))
				session_indices+=("${candidate_indices[$idx]}")
				unset "candidate_indices[$idx]"
				candidate_indices=(${candidate_indices[@]-})
			done
		fi

		if (( ${#session_indices[@]} == 0 )); then
			continue
		fi

		local session_cards_json="[]"
		for card_array_index in "${session_indices[@]}"; do
			local card_id="${deck_card_ids[$card_array_index]}"
			session_cards_json=$(jq -c --arg id "$card_id" '. + [$id]' <<<"$session_cards_json")
		done

		create_study_session "$deck_id" "$session_id" "$session_cards_json" "${#session_indices[@]}" "completed"

		local current_ts="$session_ts"
		for card_array_index in "${session_indices[@]}"; do
			local card_id="${deck_card_ids[$card_array_index]}"
			local level="${deck_card_levels[$card_array_index]:-1}"
			local difficulty="${deck_card_difficulty[$card_array_index]:-3}"
			local answer_type="${deck_card_answer_type[$card_array_index]:-multiple_choice}"
			local correct_answer="${deck_card_correct[$card_array_index]:-}"
			local wrong_options="${deck_card_wrong_options[$card_array_index]:-[]}"
			local chance roll result user_answer self_assessed answered_at answered_ts interval_days

			chance=$(calc_correct_chance "$level" "$difficulty" "$answer_type")
			roll=$((RANDOM % 100))

			if (( roll < chance )); then
				result="correct"
				user_answer="$correct_answer"
			else
				result="incorrect"
				if [[ "$answer_type" == "short_answer" ]]; then
					user_answer=$(pick_short_answer_wrong)
				else
					user_answer=$(pick_wrong_answer "$wrong_options" "${deck_card_wrong[$card_array_index]:-nao sei}")
				fi
			fi

			if [[ "$answer_type" == "short_answer" ]]; then
				self_assessed=true
			else
				self_assessed=false
			fi

			current_ts=$((current_ts + (60 + (RANDOM % 120))))
			answered_ts="$current_ts"
			answered_at=$(format_epoch "$answered_ts")

			create_answer_event "$session_id" "$card_id" "$result" "$user_answer" "$self_assessed" "$answered_at"

			deck_card_review_count[$card_array_index]=$((deck_card_review_count[$card_array_index] + 1))
			if [[ "$result" == "correct" ]]; then
				if (( level < 7 )); then
					level=$((level + 1))
				fi
			else
				level=$((level - 1))
				if (( level < 1 )); then
					level=1
				fi
			fi

			interval_days=$(leitner_interval_days "$level")
			deck_card_levels[$card_array_index]="$level"
			deck_card_last_review_ts[$card_array_index]="$answered_ts"
			deck_card_next_review_ts[$card_array_index]=$((answered_ts + (interval_days * 86400)))

			increment_review_day "$(format_date_only "$answered_ts")"
		done
	done

	while (( snapshot_index < ${#snapshot_ts_list[@]} )); do
		local snapshot_ts="${snapshot_ts_list[$snapshot_index]}"
		create_snapshot_from_state "$deck_id" "$snapshot_ts"
		snapshot_index=$((snapshot_index + 1))
	done

	for i in "${!deck_card_ids[@]}"; do
		local card_id="${deck_card_ids[$i]}"
		local level="${deck_card_levels[$i]:-1}"
		local last_review_ts="${deck_card_last_review_ts[$i]:-0}"
		local next_review_ts="${deck_card_next_review_ts[$i]:-0}"
		local last_reviewed_at next_review_at

		if (( last_review_ts == 0 )); then
			last_review_ts="$deck_start_ts"
		fi
		if (( next_review_ts < last_review_ts )); then
			next_review_ts=$((last_review_ts + 86400))
		fi

		last_reviewed_at=$(format_epoch "$last_review_ts")
		next_review_at=$(format_epoch "$next_review_ts")

		create_leitner_box "$deck_id" "$card_id" "$level" "$last_reviewed_at" "$next_review_at"
	done
}


echo "Seeding flashcards via $API_URL as $EMAIL (suffix: $SEED_SUFFIX)"

# Create categories first
category_count=$(echo "$CATEGORIES_JSON" | jq 'length')
if (( category_count > 0 )); then
	echo "Creating $category_count categories..."
	while IFS= read -r category; do
		category_name=$(echo "$category" | jq -r '.name')
		category_id=$(create_category "$category")
		category_names+=("$category_name")
		category_ids+=("$category_id")
		echo "  Created category: $category_name (id: $category_id)"
	done < <(echo "$CATEGORIES_JSON" | jq -c '.[]')
fi

index=0
while IFS= read -r deck; do
	index=$((index + 1))
	deck_result=$(create_deck "$deck")
	deck_id=${deck_result%%|*}
	rest=${deck_result#*|}
	deck_title=${rest%%|*}
	deck_slug=${deck_result##*|}
	deck_card_ids=()
	deck_card_correct=()
	deck_card_wrong=()
	deck_card_answer_type=()
	deck_card_levels=()
	deck_card_difficulty=()
	deck_card_wrong_options=()
	deck_card_last_review_ts=()
	deck_card_next_review_ts=()
	deck_card_review_count=()

	echo "Created deck $index: $deck_title (slug: $deck_slug, id: $deck_id)"

	# Link deck to categories
	while IFS= read -r category_name; do
		if [[ -n "$category_name" ]]; then
			cat_id=$(get_category_id_by_name "$category_name")
			if [[ -n "$cat_id" ]]; then
				link_deck_to_category "$deck_id" "$cat_id"
				echo "  Linked to category: $category_name"
			fi
		fi
	done < <(echo "$deck" | jq -r '.categories[]? // empty')

	while IFS= read -r card; do
		card_correct=$(echo "$card" | jq -r '
			if (.correct_answer | type) == "string" and (.correct_answer | length) > 0 then .correct_answer
			elif (.answers | type) == "array" and (.answers | length) > 0 then .answers[0]
			else ""
			end')
		card_wrong=$(echo "$card" | jq -r --arg correct "$card_correct" '
			if (.answers | type) == "array" then
				([.answers[] | select(. != $correct)] | .[0] // "")
			else
				""
			end')
		card_answer_type=$(echo "$card" | jq -r '.answer_type // "multiple_choice"')
		card_difficulty=$(echo "$card" | jq -r '.difficulty_level // 3')
		card_wrong_options=$(echo "$card" | jq -c --arg correct "$card_correct" '
			if (.answers | type) == "array" then
				[.answers[] | select(. != $correct)]
			else
				[]
			end')

		card_id=$(create_card "$card")
		deck_card_ids+=("$card_id")
		deck_card_correct+=("$card_correct")
		deck_card_wrong+=("${card_wrong:-nao sei}")
		deck_card_answer_type+=("$card_answer_type")
		deck_card_difficulty+=("$card_difficulty")
		deck_card_wrong_options+=("$card_wrong_options")
		link_card_to_deck "$card_id" "$deck_id"
		echo "  Linked card: $card_id"
	done < <(echo "$deck" | jq -c '.cards[]')

	seed_deck_usage_generated "$deck_id" "$deck_title"
done < <(echo "$DECKS_JSON" | jq -c '.[]')

echo "Done."
