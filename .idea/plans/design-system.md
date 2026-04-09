# Plan for creating a design system

## Rules

### Typography

- Font sizes use a 1.333 (Perfect Fourth) scale ratio
- Base size: 12px
- Scale: 9px, 12px, 16px, 21px, 28px, 37px, 48px, 64px
- 9px reserved for supplementary text (units, metadata)
- 12px for body text
- Fluid sizing approach: Use `clamp()` with viewport units for responsive scaling
- Never assume `1rem == 16px` - respect user font preferences

### Fluid Design System

- Design width reference: 600px (adjust based on project needs)
- Design height reference: 400px (optional, for full-screen interfaces)
- Width ratio (`--vw-ratio`): Controls viewport width influence (0-1)
- Height ratio (`--vh-ratio`): Controls viewport height influence (0-1)
- Fluid calculation creates "relative pixel" values that scale with viewport
- Divide fluid calculation by 16 to get pseudo-px factor for easy integration

**Fluid CSS Variables:**

```css
:root {
  --design-width: 600;
  --design-height: 400;
  --vw-ratio: 0.8;
  --vh-ratio: 0.8;
}

body {
  --fluid: calc([fluid-calculation] / 16);
}

/* Usage example */
.element {
  width: calc(60 * var(--fluid));
  padding: calc(16 * var(--fluid));
  font-size: calc(14 * var(--fluid));
}
```

### Color Schema

- 5 colors total: 2 neutrals, 1 base, 1 contrast, 1 accent
- Use OKLCH color space for perceptually uniform colors
- Always provide fallback values for browsers without OKLCH support
- OKLCH benefits:
  - Perceptually uniform color changes
  - Consistent lightness across hues
  - Better accessibility (lightness ties to contrast)
  - Independent manipulation (lightness, chroma, hue)
  - Wide gamut support (P3 color space)

**Color implementation pattern:**

```css
:root {
  /* Fallback in hex/rgb */
  --color-base: #3b82f6;

  /* OKLCH with feature detection */
  @supports (color: oklch(0 0 0)) {
    --color-base: oklch(0.65 0.18 250);
  }
}
```

**Color structure (to be defined):**

- **Neutral**: Light backgrounds, borders
- **Surface**: Dark text, high contrast elements
- **Base**: Primary brand color, main actions
- **Contrast**: Pairs with base for emphasis
- **Accent**: Highlights, secondary actions

**Extended palette (for data visualization):**

- Maintain consistent lightness levels for accessibility
- Use chroma variations for data series differentiation

### Spacing System

- 8px base unit for visual hierarchy
- Apply fluid scaling to spacing: `calc(8 * var(--fluid))`, `calc(16 * var(--fluid))`, etc.

### Shadows & Elevation

Use layered box-shadows for smoother, more realistic depth. Technique: multiply offset and blur by 2 for each layer.

**Shadow levels:**

**Level 1 - Subtle (hover states, slight lift):**

```css
box-shadow:
  0 1px 1px rgba(0,0,0,0.11),
  0 2px 2px rgba(0,0,0,0.11),
  0 4px 4px rgba(0,0,0,0.11),
  0 8px 8px rgba(0,0,0,0.11);
```

**Level 2 - Cards (default elevation):**

```css
box-shadow:
  0 1px 1px rgba(0,0,0,0.12),
  0 2px 2px rgba(0,0,0,0.12),
  0 4px 4px rgba(0,0,0,0.12),
  0 8px 8px rgba(0,0,0,0.12),
  0 16px 16px rgba(0,0,0,0.12);
```

**Level 3 - Popovers (sharp, concentrated):**

```css
box-shadow:
  0 1px 1px rgba(0,0,0,0.25),
  0 2px 2px rgba(0,0,0,0.20),
  0 4px 4px rgba(0,0,0,0.15),
  0 8px 8px rgba(0,0,0,0.10),
  0 16px 16px rgba(0,0,0,0.05);
```

**Level 4 - Modals (diffuse, distant):**

```css
box-shadow:
  0 1px 2px rgba(0,0,0,0.07),
  0 2px 4px rgba(0,0,0,0.07),
  0 4px 8px rgba(0,0,0,0.07),
  0 8px 16px rgba(0,0,0,0.07),
  0 16px 32px rgba(0,0,0,0.07),
  0 32px 64px rgba(0,0,0,0.07);
```

**Usage guidelines:**

- More layers = softer spread (reduce alpha per layer to maintain strength)
- Decrease alpha per layer for sharp shadows (concentrated)
- Increase alpha per layer for diffuse shadows (soft)
- Increase blur increments for dreamy effects
- Decouple Y-offset from blur to control shadow distance

### Motion & Easing

- Animation duration should not exceed 200ms
- Always respect user motion preferences with `@media (prefers-reduced-motion)`
- Easing functions:
  - `--ease-fluid: cubic-bezier(0.3, 0, 0, 1)` - smooth, flowing transitions
  - `--ease-snappy: cubic-bezier(0.2, 0, 0, 1)` - quick, responsive interactions
  - `--ease-enter: cubic-bezier(0, 0.5, 0.5, 1)` - elements entering the viewport
  - `--ease-move: cubic-bezier(0, 0.5, 0.5, 1)` - elements changing position
  - `--ease-exit: cubic-bezier(0.5, 0, 0, 1)` - elements leaving the viewport

**Reduced motion implementation:**

```css
.element {
  transition: transform 200ms var(--ease-snappy);
}

@media (prefers-reduced-motion: reduce) {
  .element {
    transition-duration: 0.01ms;
    /* Or remove animations entirely */
    animation: none;
  }
}
```

**Guidelines for reduced motion:**

- Disable decorative animations completely
- Keep essential state transitions (loading, success/error feedback)
- Reduce duration to near-instant (0.01ms) instead of removing entirely
- Maintain opacity transitions as they don't trigger motion sickness
- Disable parallax, auto-playing carousels, and scrolling animations

### Transitions & Animations

Use CSS variables for all transition values to enable centralized control.

**Transition types:**

#### A. Style Transitions (opacity, color, background, shadow, border, transform)

Purpose: Notify users of changes, indicate interactivity, show state changes

**When to use:**

- Buttons (hover, active, disabled states)
- Links (hover, visited)
- Popovers, tooltips, menus (appear/disappear)
- Modals, dialogs (fade in/out)
- Cards (hover elevation)

**Implementation:**

```css
.button {
  transition:
    background-color 200ms var(--ease-snappy),
    box-shadow 200ms var(--ease-snappy),
    transform 200ms var(--ease-snappy);
}
```

**Principles:**

- Keep subtle - serve the purpose, don't distract
- Use pre-defined tokens (duration + easing variables)
- Consistent across similar elements

#### B. Layout Transitions (position, size, height, width)

Purpose: Show UI logic and relationships, build spatial understanding

**When to use:**

- Expand/collapse (file trees, accordions, show more)
- Slide in/out (drawers, sidebars, submenus)
- Modal entry/exit (scale + fade)
- Dynamic content changes

**Fixed layouts:**
Craft case-by-case with defined durations:

```css
.drawer {
  transform: translateX(-100%);
  transition: transform 250ms var(--ease-enter);
}

.drawer.open {
  transform: translateX(0);
}
```

**Dynamic layouts:**
Use spring-based animations with ResizeObserver for:

- Content with unpredictable height (collapse, file trees)
- Responsive layout changes
- Elements that resize during animation

**Challenges with CSS:**

- Cannot animate from `0` to `auto` - requires measurement
- Need to monitor layout changes during animation
- Solution: Use spring-based libraries (react-spring) + ResizeObserver

#### When NOT to animate

**Toolbar menus & submenus:**

- No fade in/out when hovering between items
- Keep visual continuity when moving cursor
- Users focus on content, not animations
- Show instantly on open, fade out on close (like macOS)

**High-frequency interactions:**

- Tooltips should appear quickly, minimal animation
- Right-click menus (frequent, low novelty)
- Table row selections
- Filter applications

**Performance concerns:**

- If DOM animations slow the site, remove them
- Consider alternatives: video, canvas, SVG
- Monitor with performance tools

**Functionality first:**

- Reduce animations on common components
- Don't add novelty where efficiency matters
- Avoid Material Design's excessive animations

#### Animation control

**Global toggle:**

```css
:root {
  --animation-duration: 200ms;
}

.reduce-motion {
  --animation-duration: 0.01ms;
}
```

Keep all animation variables centralized for easy toggling.

### Image Resolution Guidelines

| Aspect Ratio | Max         | Min       |
|--------------|-------------|-----------|
| 16:9         | 1920 x 1080 | 480 x 270 |
| 1:1          | 1080 x 1080 | 270 x 270 |
| 3:2          | 1080 x 720  | 270 x 180 |
| 4:3          | 2048 x 1536 | 512 x 512 |

### Specific Image Sizes

| Use Case                    | Dimensions    | Aspect Ratio |
|-----------------------------|---------------|--------------|
| Background Image            | 1920 x 1080   | 16:9         |
| Hero Image                  | 1280 x 720    | 16:9         |
| Website Banner              | 250 x 250     | 1:1          |
| Blog Image                  | 1200 x 630    | 3:2          |
| Logo (Rectangle)            | 250 x 100     | 2:3          |
| Logo (Square)               | 100 x 100     | 1:1          |
| Favicon                     | 16 x 16       | 1:1          |
| Social Media Icons          | 32 x 32       | 1:1          |
| Lightbox Images (Full Screen) | 1600 x 500  | 16:9         |
| Thumbnail Image             | 150 x 150     | 1:1          |

## Reference Guidelines

### Typography (Reference)

- Apply `-webkit-font-smoothing: antialiased` for better legibility
- Apply `text-rendering: optimizeLegibility`
- Font weights below 400 should not be used
- Use `font-variant-numeric: tabular-nums` for tables and timers
- Prevent text resizing on iOS landscape with `-webkit-text-size-adjust: 100%`
- Font weight should not change on hover to prevent layout shift
- Test fluid calculations with browser zoom and user font preferences

### Interactivity

- Clicking input labels should focus the input field
- Wrap inputs with `<form>` to submit by pressing Enter
- Use appropriate input types (`password`, `email`, etc.)
- Disable buttons after submission to avoid duplicate requests
- Interactive elements should disable `user-select` for inner content
- Decorative elements should disable `pointer-events`

### Motion

- Animation duration should not exceed 200ms for immediate feel
- Switching themes should not trigger transitions on elements
- Frequent actions should avoid extraneous animations
- Use `scroll-behavior: smooth` for in-page anchor navigation
- Always respect `prefers-reduced-motion` user preference

### Touch

- Hover states should use `@media (hover: hover)` to avoid touch press visibility
- Input font size should be minimum 16px to prevent iOS zoom on focus
- Inputs should not auto-focus on touch devices
- Apply `muted` and `playsinline` to `<video>` for iOS autoplay
- Disable iOS tap highlight with `-webkit-tap-highlight-color: rgba(0,0,0,0)` and provide alternative

### Optimizations

- Large `blur()` values may be slow
- Sparingly enable GPU rendering with `transform: translateZ(0)`
- Toggle `will-change` only for duration of animations
- Pause off-screen looping animations to reduce CPU/GPU usage

### Accessibility

- Use box shadow for focus rings, not outline
- Icon-only elements should have explicit `aria-label`
- Images should use `<img>` for screen readers
- Gradient text should unset gradient on `::selection` state
- Disabled buttons should not have tooltips
- Test fluid calculations at fluid.style to ensure zoom accessibility
- Always respect `prefers-reduced-motion` for accessibility

### Miscelanious

- Style document selection with `::selection`
- Display feedback relative to its trigger
- Empty states should prompt to create new items
