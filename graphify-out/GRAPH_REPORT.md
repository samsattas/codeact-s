# Graph Report - CodeAct-Style  (2026-08-10)

## Corpus Check
- 125 files · ~110,582 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2218 nodes · 2761 edges · 145 communities (120 shown, 25 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 85 edges (avg confidence: 0.69)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f134a6c6`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- gray
- BM25
- Sandbox
- BM25
- slide_search_core.py
- color
- spacing
- html-token-validator.py
- TestTailwindConfigGenerator
- BM25
- Tailwind CSS Utility Reference
- newTestServer
- generate-slide.py
- DesignSystemGenerator
- TailwindConfigGenerator
- fetch-background.py
- icon/generate.py
- fontSize
- TestShadcnInstaller
- main
- test_design_system_mode.py
- Brand Guidelines v1.0
- extract-colors.cjs
- validate-asset.cjs
- design_system.py
- button
- .add_components
- ShadcnInstaller
- diagnosisHint
- validate-tokens.cjs
- test_tailwind_config_gen.py
- .generate_config_string
- inject-brand-context.cjs
- embed-tokens.cjs
- primitive
- patch
- NewOllamaProvider
- ._base_config
- logo/generate.py
- generate-tokens.cjs
- Design
- Canvas Design System
- sync-brand-to-tokens.cjs
- _run
- Form & Input Components
- Tailwind CSS Responsive Design
- radius
- Typography Specifications
- CodeAct Agent
- _select_palette_for_mode
- shadow
- Logo Usage Rules
- Component Specifications
- shadcn/ui Accessibility Patterns
- lg
- Asset Approval Checklist
- Logo AI Prompt Engineering
- xl
- newTestAnthropicProvider
- none
- validate_data.py
- test_sync_brand_to_tokens.py
- main
- Color Palette Management
- .generate
- CIP Deliverable Guide
- States and Variants
- What You Must Do When Invoked
- .test_add_fonts
- UI Styling Skill
- Workflow
- .test_recommend_plugins
- Design System
- Tailwind CSS Customization
- .test_write_config_creates_content
- Routing by Task Type
- shadcn/ui Theming & Customization
- Asset Organization Guide
- codeact-agent
- Primary Color Meanings
- .test_add_components_no_components
- Core Logo Types
- Brand Consistency Checklist
- CIP Mockup Prompt Engineering
- Color Semantics
- Design Principles
- Design Principles
- CIP Design Reference
- Icon Design Reference
- Copywriting Formulas
- Copywriting Formulas
- Banner Design - Multi-Format Creative Banner System
- Messaging Framework
- Brand Voice Framework
- Layout Patterns
- Tailwind Integration
- Layout Patterns
- search
- brand/references/update.md
- Logo Design Reference
- Token Architecture
- Primitive Tokens
- Core Visual Elements
- CIP Design Style Guide
- Quick Reference
- Brand
- Slide Strategies
- Component Tokens
- Slide Strategies
- graphify reference: extra exports and benchmark
- Prerequisites
- UI/UX Pro Max - Design Intelligence
- Slides Reference
- HTML Slide Template
- HTML Slide Template
- detect_domain
- How to Use This Skill
- graphify reference: query, path, explain
- Slides
- Pre-Delivery Checklist
- Brand Guidelines Template
- Common Rules for Professional UI
- Example Workflow
- graphify reference: add a URL and watch a folder
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- graphify reference: GitHub clone and cross-repo merge
- graphify reference: transcribe video and audio
- .__init__
- graphify
- slides-create.md
- extraction-spec.md
- create.md
- .test_add_components_no_config
- .test_add_components_dry_run
- .test_list_installed_no_config
- .test_init_dry_run
- .test_generate_typescript_config
- .test_generate_config_with_plugins
- .test_init_javascript
- .test_write_config_invalid_path
- .test_default_output_path_typescript
- .test_default_content_paths_vue
- .test_add_colors
- md

## God Nodes (most connected - your core abstractions)
1. `TailwindConfigGenerator` - 58 edges
2. `TestTailwindConfigGenerator` - 35 edges
3. `ShadcnInstaller` - 34 edges
4. `DesignSystemGenerator` - 29 edges
5. `TestShadcnInstaller` - 26 edges
6. `Sandbox` - 18 edges
7. `UI Styling Skill` - 17 edges
8. `color` - 15 edges
9. `search()` - 15 edges
10. `Design` - 15 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `FromEnv()`  [INFERRED]
  cmd/agent/main.go → internal/llm/llm.go
- `main()` --calls--> `NewSandbox()`  [INFERRED]
  cmd/web/main.go → internal/tools/tools.go
- `newTestServer()` --calls--> `NewSandbox()`  [INFERRED]
  cmd/web/main_test.go → internal/tools/tools.go
- `main()` --calls--> `NewSandbox()`  [INFERRED]
  cmd/agent/main.go → internal/tools/tools.go
- `main()` --calls--> `FromEnv()`  [INFERRED]
  cmd/web/main.go → internal/llm/llm.go

## Import Cycles
- None detected.

## Communities (145 total, 25 thin omitted)

### Community 0 - "gray"
Cohesion: 0.05
Nodes (53): $type, $value, $type, $value, $type, $value, $type, $value (+45 more)

### Community 1 - "BM25"
Cohesion: 0.07
Nodes (42): BM25, detect_domain(), get_cip_brief(), _load_csv(), Load CSV and return list of dicts, Core search function using BM25, Auto-detect the most relevant domain from query, Main search function with auto-domain detection (+34 more)

### Community 2 - "Sandbox"
Cohesion: 0.06
Nodes (54): Agent, fakeProvider, Outcome, Step, fatal(), main(), repl(), runTask() (+46 more)

### Community 3 - "BM25"
Cohesion: 0.13
Nodes (16): BM25, _domain_keywords(), _get_bm25(), _load_csv(), _load_product_keywords(), _normalize(), Apply synonym substitution before tokenizing., BM25 ranking algorithm for text search (+8 more)

### Community 4 - "slide_search_core.py"
Cohesion: 0.09
Nodes (36): format_context(), format_result(), main(), Format a single search result for display, Format contextual recommendations for display., BM25, calculate_pattern_break(), detect_domain() (+28 more)

### Community 5 - "color"
Cohesion: 0.04
Nodes (48): $type, $value, background, destructive, destructive-foreground, foreground, muted, muted-foreground (+40 more)

### Community 6 - "spacing"
Cohesion: 0.06
Nodes (34): $type, $value, $type, $value, $type, $value, $type, $value (+26 more)

### Community 7 - "html-token-validator.py"
Cohesion: 0.14
Nodes (24): get_context(), is_allowed_exception(), is_allowed_rgba(), is_inside_block(), load_css_variables(), main(), print_result(), print_summary() (+16 more)

### Community 8 - "TestTailwindConfigGenerator"
Cohesion: 0.07
Nodes (15): Test adding colors multiple times., Test adding full color palette., Test adding custom breakpoints., Test TailwindConfigGenerator class., Test that adding same plugin twice doesn't duplicate., Test plugin recommendations for Next.js., Test initialization with default settings., Test generating JavaScript configuration. (+7 more)

### Community 9 - "BM25"
Cohesion: 0.12
Nodes (19): BM25, detect_domain(), _load_csv(), Load CSV and return list of dicts, Core search function using BM25, Auto-detect the most relevant domain from query, Main search function with auto-domain detection, Search across all domains and combine results (+11 more)

### Community 10 - "Tailwind CSS Utility Reference"
Cohesion: 0.05
Nodes (43): Arbitrary Values, Aspect Ratio, Background Colors, Border Color, Border Radius, Border Style, Border Width, Borders (+35 more)

### Community 11 - "newTestServer"
Cohesion: 0.11
Nodes (26): main(), Context, T, newTestServer(), TestHandleBrowseListsSubdirectoriesOnly(), TestHandleBrowseRejectsMissingPath(), TestHandleInfo(), TestHandleRunRejectsInvalidDir() (+18 more)

### Community 12 - "generate-slide.py"
Cohesion: 0.15
Nodes (19): _e(), generate_chart_slide(), generate_cta_slide(), generate_deck(), generate_metrics_slide(), generate_problem_slide(), generate_solution_slide(), generate_testimonial_slide() (+11 more)

### Community 13 - "DesignSystemGenerator"
Cohesion: 0.14
Nodes (11): DesignSystemGenerator, _palette_is_dark(), WCAG relative luminance of a #RRGGBB string, or None if unparseable., True when a colors.csv row's Background is a dark surface., Generates design system recommendations from aggregated searches., Load reasoning rules from CSV., _relative_luminance(), TestReasoningMatch (+3 more)

### Community 14 - "TailwindConfigGenerator"
Cohesion: 0.10
Nodes (11): Generate Tailwind CSS configuration files., Add full color palette (50-950 shades) for a base color. Args: name: Color name…, TailwindConfigGenerator, Test adding custom spacing., Test validating config with no content paths., Test validating config with empty theme extensions., Test writing configuration to file., Test initialization with different frameworks. (+3 more)

### Community 15 - "fetch-background.py"
Cohesion: 0.17
Nodes (17): generate_css_for_background(), get_background_image(), get_curated_images(), get_overlay_css(), get_pexels_search_url(), load_backgrounds_config(), load_brand_colors(), main() (+9 more)

### Community 16 - "icon/generate.py"
Cohesion: 0.20
Nodes (15): apply_color(), apply_viewbox_size(), extract_svgs(), generate_batch(), generate_icon(), generate_sizes(), load_env(), main() (+7 more)

### Community 17 - "fontSize"
Cohesion: 0.12
Nodes (16): $type, $value, $type, $value, $type, $value, $type, $value (+8 more)

### Community 18 - "TestShadcnInstaller"
Cohesion: 0.12
Nodes (10): Test ShadcnInstaller class., Test adding all components without config., Test adding all components in dry run mode., Create temporary project structure., Test listing installed components when none exist., Test listing installed components when they exist., Test checking for existing shadcn config., Test getting installed components without config. (+2 more)

### Community 19 - "main"
Cohesion: 0.13
Nodes (8): main(), Add custom font families. Args: fonts: Dict of font_type: [font_names] e.g.,…, Add custom spacing values. Args: spacing: Dict of name: value e.g., {'18':…, Add custom breakpoints. Args: breakpoints: Dict of name: width e.g., {'3xl':…, Add plugin requirements. Args: plugins: List of plugin names e.g.,…, Get plugin recommendations based on configuration. Returns: List of recommended…, Validate configuration. Returns: Tuple of (valid, message), Add custom colors to theme. Args: colors: Dict of color_name: color_value Value…

### Community 20 - "test_design_system_mode.py"
Cohesion: 0.16
Nodes (10): _filter_anti_patterns_for_mode(), _query_wants_dark(), True when a styles.csv row describes itself as dark-first., True when the query explicitly asks for a dark theme., Resolve the mode the rest of the output has to agree with., Drop "avoid dark mode" advice once dark mode is the resolved answer., _resolve_color_mode(), _style_is_dark_primary() (+2 more)

### Community 21 - "Brand Guidelines v1.0"
Cohesion: 0.05
Nodes (37): 1. Color Palette, 2. Typography, 3. Logo Usage, 4. Voice & Tone, 5. Imagery Guidelines, 6. Design Components, Accessibility, AI Image Generation (+29 more)

### Community 22 - "extract-colors.cjs"
Cohesion: 0.22
Nodes (11): calculateCompliance(), colorDistance(), displayPalette(), extractHexColors(), findNearestBrandColor(), fs, generateImageMagickCommand(), hexToRgb() (+3 more)

### Community 23 - "validate-asset.cjs"
Cohesion: 0.25
Nodes (13): checkManifest(), formatBytes(), formatOutput(), fs, main(), parseFilename(), path, RULES (+5 more)

### Community 24 - "design_system.py"
Cohesion: 0.11
Nodes (25): ansi_ljust(), _detect_page_type(), format_ascii_box(), format_markdown(), format_master_md(), format_page_override_md(), generate_design_system(), _generate_intelligent_overrides() (+17 more)

### Community 25 - "button"
Cohesion: 0.06
Nodes (45): $type, $value, $type, $value, bg, fg, font-size, hover-bg (+37 more)

### Community 26 - ".add_components"
Cohesion: 0.17
Nodes (8): main(), Add all available shadcn/ui components. Args: overwrite: If True, overwrite…, List installed components. Returns: Tuple of (success, message with component…, Check if shadcn is initialized in project. Returns: True if components.json…, Get list of already installed components. Returns: List of installed component…, Read shadcn version from project package.json; fall back to a pinned default., Add shadcn/ui components. Args: components: List of component names to add…, Tests for shadcn_add.py

### Community 27 - "ShadcnInstaller"
Cohesion: 0.14
Nodes (8): Handle shadcn/ui component installation., ShadcnInstaller, Test adding components that are already installed., Test initialization with default project root., Test initialization with custom project root., Test checking for non-existent shadcn config., Test getting installed components when none exist., Test getting installed components when files exist.

### Community 28 - "diagnosisHint"
Cohesion: 0.31
Nodes (11): diagnosisHint(), fenced(), retryPrompt(), T, TestDiagnosisHintExitStatus(), TestDiagnosisHintUndefined(), TestDiagnosisHintUnknownErrorIsEmpty(), TestDiagnosisHintUnusedImport() (+3 more)

### Community 29 - "validate-tokens.cjs"
Cohesion: 0.24
Nodes (11): extensions, formatReport(), fs, getFiles(), main(), parseArgs(), path, patterns (+3 more)

### Community 30 - "test_tailwind_config_gen.py"
Cohesion: 0.20
Nodes (8): Tests for tailwind_config_gen.py, Reduce a generated TS/JS config to a bare assignable object so it can be handed…, Regression guard for the missing-comma bug between the ``theme`` block and…, The property preceding ``plugins`` must end with a comma (pure-Python check, so…, The emitted config parses as valid JS via ``node --check``., _strip_to_object(), TestGeneratedConfigIsValidJs, parametrize

### Community 31 - ".generate_config_string"
Cohesion: 0.20
Nodes (6): Generate configuration file content. Returns: Configuration file as string, Generate TypeScript configuration., Generate JavaScript configuration., Format plugins array for config. Validates each plugin name against a strict…, Add indentation to JSON string., Write configuration to file. Returns: Tuple of (success, message)

### Community 32 - "inject-brand-context.cjs"
Cohesion: 0.31
Nodes (10): extractColorsFromTable(), extractCoreAttributes(), extractHexColors(), extractImageStyle(), extractTypography(), extractVoice(), fs, generatePromptAddition() (+2 more)

### Community 33 - "embed-tokens.cjs"
Cohesion: 0.18
Nodes (8): args, fs, minimal, MINIMAL_TOKENS, path, projectRoot, tokensPath, wrapStyle

### Community 34 - "primitive"
Cohesion: 0.18
Nodes (11): fast, normal, slow, $type, $value, $type, $value, primitive (+3 more)

### Community 35 - "patch"
Cohesion: 0.18
Nodes (6): Test adding components with overwrite flag., Test successful component addition., Test component addition with subprocess error., Test component addition when npx is not found., Test successful addition of all components., patch

### Community 36 - "NewOllamaProvider"
Cohesion: 0.19
Nodes (14): Client, Context, NewOllamaProvider(), T, TestOllamaProviderGenerateAPIError(), TestOllamaProviderGenerateHTTPError(), TestOllamaProviderGenerateSuccess(), TestOllamaProviderName() (+6 more)

### Community 37 - "._base_config"
Cohesion: 0.22
Nodes (6): Any, Path, Initialize generator. Args: typescript: If True, generate .ts config, else .js…, Determine default output path., Create base configuration structure., Get default content paths for framework.

### Community 38 - "logo/generate.py"
Cohesion: 0.29
Nodes (9): enhance_prompt(), generate_batch(), generate_logo(), load_env(), main(), Enhance the logo prompt with style and industry modifiers, Generate a logo using Gemini models with image generation Args: aspect_ratio:…, Generate multiple logo variants with different styles (+1 more)

### Community 39 - "generate-tokens.cjs"
Cohesion: 0.36
Nodes (9): flattenTokens(), fs, generateCSS(), generateTailwind(), main(), parseArgs(), path, resolveReference() (+1 more)

### Community 40 - "Design"
Cohesion: 0.06
Nodes (35): Banner Design (Built-in), Banner: Design Rules, Banner: Quick Size Reference, Banner: Top Art Styles, Banner: Workflow, CIP Design (Built-in), CIP: Generate Brief, CIP: Generate Mockups (+27 more)

### Community 41 - "Canvas Design System"
Cohesion: 0.06
Nodes (35): 1. Visual Communication First, 2. Minimal Text Integration, 3. Expert Craftsmanship, 4. Systematic Patterns, Analog Meditation, Approach, Canvas Boundaries, Canvas Design System (+27 more)

### Community 42 - "sync-brand-to-tokens.cjs"
Cohesion: 0.33
Nodes (8): adjustBrightness(), { execFileSync }, extractColorsFromMarkdown(), fs, generateColorScale(), main(), path, updateDesignTokens()

### Community 43 - "_run"
Cohesion: 0.28
Nodes (8): Path, Regression tests for validate-tokens.cjs. The validator used to skip any line…, A hardcoded hex on the same line as a var() token is still a violation., A line that references only tokens produces no false positives., _run(), test_flags_hardcoded_hex_sharing_line_with_token(), test_token_only_line_reports_no_violation(), CompletedProcess

### Community 44 - "Form & Input Components"
Cohesion: 0.06
Nodes (32): Accordion, Alert, Alert Dialog, Avatar, Badge, Button, Card, Checkbox (+24 more)

### Community 45 - "Tailwind CSS Responsive Design"
Cohesion: 0.06
Nodes (32): 1. Mobile-First Design, 2. Consistent Breakpoint Usage, 3. Test at Breakpoint Boundaries, 4. Use Container for Content Width, 5. Progressive Enhancement, 6. Avoid Too Many Breakpoints, Best Practices, Breakpoint System (+24 more)

### Community 46 - "radius"
Cohesion: 0.29
Nodes (8): $type, $value, $type, $value, radius, default, full, default

### Community 47 - "Typography Specifications"
Cohesion: 0.06
Nodes (30): Accessibility, Base System, Best Practices, Clean & Modern, Common Font Pairings, Contrast Requirements, CSS Implementation, Editorial (+22 more)

### Community 48 - "CodeAct Agent"
Cohesion: 0.07
Nodes (25): Architecture, Commands, graphify, Architecture, Design choices, Known limitations (proof of concept, not production), Module layout, Observed behavior (+17 more)

### Community 49 - "_select_palette_for_mode"
Cohesion: 0.43
Nodes (3): Pick the highest-ranked palette matching the resolved mode. Only the dark case…, _select_palette_for_mode(), TestPaletteSelection

### Community 50 - "shadow"
Cohesion: 0.47
Nodes (6): sm, shadow, sm, sm, $type, $value

### Community 51 - "Logo Usage Rules"
Cohesion: 0.07
Nodes (28): Absolute Don'ts, Approved Backgrounds, Before Using Logo, Clear Space, Co-branding, Color Rules, Color Usage, Color Variants (+20 more)

### Community 52 - "Component Specifications"
Cohesion: 0.07
Nodes (28): Alert, Anatomy, Anatomy, Anatomy, Anatomy, Anatomy, Badge, Button (+20 more)

### Community 53 - "shadcn/ui Accessibility Patterns"
Cohesion: 0.07
Nodes (28): Accordion, Alert, ARIA Labels, Checkbox and Radio, Color Contrast, Command Palette Navigation, Component-Specific Patterns, Dialog/Modal Navigation (+20 more)

### Community 54 - "lg"
Cohesion: 0.60
Nodes (5): lg, $type, $value, lg, lg

### Community 55 - "Asset Approval Checklist"
Cohesion: 0.08
Nodes (25): Accessibility, Archival, Asset Approval Checklist, Automation Support, Color Compliance, Common Issues & Fixes, Content Accessibility, Content Quality (+17 more)

### Community 56 - "Logo AI Prompt Engineering"
Cohesion: 0.08
Nodes (25): Common Pitfalls, Core Prompt Structure, Detailed Brief, Eco/Sustainable, Effective Keywords by Style, Fashion Brand, Healthcare, Industry-Specific Prompts (+17 more)

### Community 57 - "xl"
Cohesion: 0.67
Nodes (4): xl, xl, $type, $value

### Community 58 - "newTestAnthropicProvider"
Cohesion: 0.19
Nodes (14): Client, Context, NewAnthropicProvider(), T, newTestAnthropicProvider(), TestAnthropicProviderGenerateAPIError(), TestAnthropicProviderGenerateConcatenatesTextBlocks(), TestAnthropicProviderGenerateHTTPError() (+6 more)

### Community 59 - "none"
Cohesion: 0.67
Nodes (4): $type, $value, none, none

### Community 60 - "validate_data.py"
Cohesion: 0.83
Nodes (3): _check_file(), main(), _read_rows()

### Community 63 - "Color Palette Management"
Cohesion: 0.08
Nodes (24): Accessibility Requirements, Brand Compliance Validation, Checking Contrast, Color Documentation Format, Color Extraction, Color Palette Examples, Color Palette Management, Color System Structure (+16 more)

### Community 64 - ".generate"
Cohesion: 0.14
Nodes (8): Execute searches across multiple domains., Find matching reasoning rule for a category., Apply reasoning rules to search results., Select best matching result based on priority keywords., Extract results list from search result dict., Generate complete design system recommendation. variance/motion/density are…, Bucket a 1-10 dial value into its tier config. Returns None if value is None., _resolve_dial()

### Community 65 - "CIP Deliverable Guide"
Cohesion: 0.08
Nodes (24): Apparel, Business Card, Car/Sedan, CIP Deliverable Guide, Core Identity, Digital Assets, Email Signature, Envelope (+16 more)

### Community 66 - "States and Variants"
Cohesion: 0.08
Nodes (24): Accessibility, Accessibility Requirements, ARIA States, Color Contrast, Color Variants, Disabled States, Error Messages, Error States (+16 more)

### Community 67 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 69 - "UI Styling Skill"
Cohesion: 0.08
Nodes (24): Accessibility Patterns, Alternative: Tailwind-Only Setup, Best Practices, Common Patterns, Component Layer: shadcn/ui, Component Library Guide, Component + Styling Setup, Core Stack (+16 more)

### Community 70 - "Workflow"
Cohesion: 0.08
Nodes (23): Art Direction Styles (Reuse from Banner), Color & Contrast, Design Best Practices, HTML Design Rules, HTML Template Structure, Option A: Chrome Headless CLI (Recommended — zero dependencies), Option B: chrome-devtools skill, Option C: Playwright script (+15 more)

### Community 72 - "Design System"
Cohesion: 0.09
Nodes (22): Best Practices, Chart.js Integration, Command, Component Spec Pattern, Contextual Decision Flow, Decision System CSVs, Design System, Integration (+14 more)

### Community 73 - "Tailwind CSS Customization"
Cohesion: 0.09
Nodes (22): @apply Directive, Best Practices, Color Customization, Complete Tailwind Config, Configuration Examples, Content Configuration, Custom Color Palette, Custom Font Sizes (+14 more)

### Community 75 - "Routing by Task Type"
Cohesion: 0.10
Nodes (19): Banner Design Tasks, Brand Identity Tasks, Component Creation, Corporate Identity Program Tasks, Design Routing Guide, Design System Migration, Icon Design Tasks, Implementation Tasks (+11 more)

### Community 76 - "shadcn/ui Theming & Customization"
Cohesion: 0.10
Nodes (19): Base Color Presets, Best Practices, Color Customization, Color Format, Component Customization, CSS Variable System, Customize Styles, Customize Variants (+11 more)

### Community 77 - "Asset Organization Guide"
Cohesion: 0.11
Nodes (18): Asset Entry (manifest.json), Asset Organization Guide, By Campaign, By Status, By Type, Cleanup Workflow, Components, Directory Structure (+10 more)

### Community 79 - "Primary Color Meanings"
Cohesion: 0.11
Nodes (18): Accessibility Considerations, Analogous, Black, Blue, Color Combinations by Industry, Color Harmony Types, Complementary, Green (+10 more)

### Community 81 - "Core Logo Types"
Cohesion: 0.11
Nodes (18): 1. Wordmark (Logotype), 2. Lettermark (Monogram), 3. Pictorial Mark (Brand Mark), 4. Abstract Mark, 5. Mascot, 6. Emblem, 7. Combination Mark, Aesthetic Styles (+10 more)

### Community 82 - "Brand Consistency Checklist"
Cohesion: 0.11
Nodes (17): Audit Frequency, Brand Consistency Checklist, Channel Audit, Collateral, Colors, Common Issues, Email, Imagery (+9 more)

### Community 83 - "CIP Mockup Prompt Engineering"
Cohesion: 0.11
Nodes (17): Apparel (Polo/T-Shirt), Base Prompt Structure, Business Card, CIP Mockup Prompt Engineering, Context Modifiers, Corporate Minimal, Deliverable-Specific Modifiers, Letterhead (+9 more)

### Community 84 - "Color Semantics"
Cohesion: 0.11
Nodes (17): Accent, Applying Semantic Tokens, Background & Foreground, Border & Ring, Color Semantics, Dark Mode Overrides, Destructive, Interactive States (+9 more)

### Community 85 - "Design Principles"
Cohesion: 0.12
Nodes (15): 22 Art Direction Styles, Banner Sizes & Art Direction Styles Reference, Complete Banner Sizes, CTA Rules, Design Principles, Pinterest Research Queries, Print, Print Specs (+7 more)

### Community 86 - "Design Principles"
Cohesion: 0.12
Nodes (15): 22 Art Direction Styles, Banner Sizes & Art Direction Styles Reference, Complete Banner Sizes, CTA Rules, Design Principles, Pinterest Research Queries, Print, Print Specs (+7 more)

### Community 87 - "CIP Design Reference"
Cohesion: 0.13
Nodes (14): CIP Brief (Start Here), CIP Design Reference, Commands, Deliverable Categories, Design Styles, Detailed References, Generate Mockups, HTML Presentation Features (+6 more)

### Community 88 - "Icon Design Reference"
Cohesion: 0.13
Nodes (14): Available Styles, CLI Options, Commands, Generate Batch Variations, Generate Multiple Sizes, Generate Single Icon, Icon Categories, Icon Design Reference (+6 more)

### Community 89 - "Copywriting Formulas"
Cohesion: 0.13
Nodes (14): AIDA (Attention-Interest-Desire-Action), Before-After-Bridge, Contrast Patterns, Copywriting Formulas, Core Formulas, Cost of Inaction, FAB (Features-Advantages-Benefits), Formula-to-Slide Mapping (+6 more)

### Community 90 - "Copywriting Formulas"
Cohesion: 0.13
Nodes (14): AIDA (Attention-Interest-Desire-Action), Before-After-Bridge, Contrast Patterns, Copywriting Formulas, Core Formulas, Cost of Inaction, FAB (Features-Advantages-Benefits), Formula-to-Slide Mapping (+6 more)

### Community 91 - "Banner Design - Multi-Format Creative Banner System"
Cohesion: 0.14
Nodes (13): Art Direction Styles (Top 10), Banner Design - Multi-Format Creative Banner System, Banner Size Quick Reference, Design Rules, Prerequisites, Security, Step 1: Gather Requirements (AskUserQuestion), Step 2: Research & Art Direction (+5 more)

### Community 92 - "Messaging Framework"
Cohesion: 0.14
Nodes (13): Core Statements, Elevator Pitches, Framework Structure, Message Architecture, Message by Audience, Message Testing, Messaging Framework, Mission Statement (+5 more)

### Community 93 - "Brand Voice Framework"
Cohesion: 0.14
Nodes (13): Brand Voice Framework, Character Spectrum, Emotion Spectrum, Language Spectrum, Step 1: Define Personality Traits, Step 2: Create Voice Chart, Step 3: Context Adaptation, Tone Spectrum (+5 more)

### Community 94 - "Layout Patterns"
Cohesion: 0.14
Nodes (13): Card Styles, Component Variants, CSS Structures, Feature Grid (3 columns), Layout Decision Flow, Layout Patterns, Layout Selection by Use Case, Metric Styles (+5 more)

### Community 95 - "Tailwind Integration"
Cohesion: 0.14
Nodes (13): Animation Tokens, Base Layer, Button Example, Component Classes, CSS Variables Setup, Dark Mode Toggle, HSL Format Benefits, shadcn/ui Alignment (+5 more)

### Community 96 - "Layout Patterns"
Cohesion: 0.14
Nodes (13): Card Styles, Component Variants, CSS Structures, Feature Grid (3 columns), Layout Decision Flow, Layout Patterns, Layout Selection by Use Case, Metric Styles (+5 more)

### Community 97 - "search"
Cohesion: 0.14
Nodes (11): All indexed terms, for suggestion/typo-recovery purposes., Nearest known vocabulary terms for a query that returned 0 hits, so the caller…, Main search function with auto-domain detection, Search stack-specific guidelines, search(), search_stack(), _suggest_terms(), format_output() (+3 more)

### Community 98 - "brand/references/update.md"
Cohesion: 0.15
Nodes (12): Color Presets, Examples, Files Modified, Important, Overview, Skills Used, Step 1: Gather Brand Input, Step 2: Update Brand Guidelines (+4 more)

### Community 99 - "Logo Design Reference"
Cohesion: 0.15
Nodes (12): Available Styles, Color Psychology, Commands, Design Brief (Start Here), Detailed References, Generate Logo, Industry Defaults, Logo Design Reference (+4 more)

### Community 100 - "Token Architecture"
Cohesion: 0.15
Nodes (12): Categories, Dark Mode, File Organization, Layer 1: Primitive Tokens, Layer 2: Semantic Tokens, Layer 3: Component Tokens, Layer Overview, Migration from Flat Tokens (+4 more)

### Community 101 - "Primitive Tokens"
Cohesion: 0.17
Nodes (11): Border Radius, Color Scales, Gray Scale, Motion / Duration, Primary Colors (Blue), Primitive Tokens, Shadows, Spacing Scale (+3 more)

### Community 102 - "Core Visual Elements"
Cohesion: 0.18
Nodes (10): Color Palette, Colors, Core Visual Elements, Logo, Logo, Quick Checks, Typography, Typography (+2 more)

### Community 103 - "CIP Design Style Guide"
Cohesion: 0.18
Nodes (10): Bold Dynamic, CIP Design Style Guide, Classic Traditional, Color Psychology, Corporate Minimal, Fresh Modern, Luxury Premium, Modern Tech (+2 more)

### Community 104 - "Quick Reference"
Cohesion: 0.18
Nodes (11): 10. Charts & Data (LOW), 1. Accessibility (CRITICAL), 2. Touch & Interaction (CRITICAL), 3. Performance (HIGH), 4. Style Selection (HIGH), 5. Layout & Responsive (HIGH), 6. Typography & Color (MEDIUM), 7. Animation (MEDIUM) (+3 more)

### Community 105 - "Brand"
Cohesion: 0.20
Nodes (9): Brand, Brand Sync Workflow, Quick Start, References, Routing, Scripts, Subcommands, Templates (+1 more)

### Community 106 - "Slide Strategies"
Cohesion: 0.20
Nodes (9): Common Structures, Duarte Sparkline Pattern, Matching Strategy to Context, Product Demo (6 slides), Sales Pitch (9 slides), Search Commands, Slide Strategies, Strategy Selection (+1 more)

### Community 107 - "Component Tokens"
Cohesion: 0.20
Nodes (9): Alert Tokens, Badge Tokens, Button Tokens, Card Tokens, Component Tokens, Dialog/Modal Tokens, Input Tokens, Table Tokens (+1 more)

### Community 108 - "Slide Strategies"
Cohesion: 0.20
Nodes (9): Common Structures, Duarte Sparkline Pattern, Matching Strategy to Context, Product Demo (6 slides), Sales Pitch (9 slides), Search Commands, Slide Strategies, Strategy Selection (+1 more)

### Community 109 - "graphify reference: extra exports and benchmark"
Cohesion: 0.22
Nodes (8): graphify reference: extra exports and benchmark, Step 6b - Wiki (only if --wiki flag), Step 7 - Neo4j export (only if --neo4j or --neo4j-push flag), Step 7a - FalkorDB export (only if --falkordb or --falkordb-push flag), Step 7b - SVG export (only if --svg flag), Step 7c - GraphML export (only if --graphml flag), Step 7d - MCP server (only if --mcp flag), Step 8 - Token reduction benchmark (only if total_words > 5000)

### Community 110 - "Prerequisites"
Cohesion: 0.22
Nodes (9): Available Domains, Available Stacks, Common Sticking Points, Output Formats, Pre-Delivery Checklist, Prerequisites, Query Strategy, Search Reference (+1 more)

### Community 111 - "UI/UX Pro Max - Design Intelligence"
Cohesion: 0.25
Nodes (7): How to Use, Primary Use Cases, Recommended, Rule Categories by Priority, Skip, UI/UX Pro Max - Design Intelligence, When to Apply

### Community 112 - "Slides Reference"
Cohesion: 0.29
Nodes (6): Key Features, Knowledge Base, Slides Reference, Usage, When to Use, Workflow

### Community 113 - "HTML Slide Template"
Cohesion: 0.29
Nodes (6): Animation Classes, Background Images, Base Structure, Chart.js Integration, CSS Variables Reference, HTML Slide Template

### Community 114 - "HTML Slide Template"
Cohesion: 0.29
Nodes (6): Animation Classes, Background Images, Base Structure, Chart.js Integration, CSS Variables Reference, HTML Slide Template

### Community 115 - "detect_domain"
Cohesion: 0.43
Nodes (3): detect_domain(), Auto-detect the most relevant domain from query. Matches are weighted by…, TestDomainDetection

### Community 116 - "How to Use This Skill"
Cohesion: 0.29
Nodes (7): How to Use This Skill, Step 1: Analyze User Requirements, Step 2: Generate Design System (REQUIRED), Step 2b: Persist Design System (Master + Overrides Pattern), Step 2c: Design Dials (optional), Step 3: Supplement with Detailed Searches (as needed), Step 4: Stack Guidelines

### Community 117 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 118 - "Slides"
Cohesion: 0.33
Nodes (5): References (Knowledge Base), Routing, Slides, Subcommands, When to Use

### Community 119 - "Pre-Delivery Checklist"
Cohesion: 0.33
Nodes (6): Accessibility, Interaction, Layout, Light/Dark Mode, Pre-Delivery Checklist, Visual Quality

### Community 120 - "Brand Guidelines Template"
Cohesion: 0.40
Nodes (4): Brand Guidelines Template, Document Structure, Extractable Fields, Usage

### Community 121 - "Common Rules for Professional UI"
Cohesion: 0.40
Nodes (5): Common Rules for Professional UI, Icons & Visual Elements, Interaction (App), Layout & Spacing, Light/Dark Mode Contrast

### Community 122 - "Example Workflow"
Cohesion: 0.40
Nodes (5): Example Workflow, Step 1: Analyze Requirements, Step 2: Generate Design System (REQUIRED), Step 3: Supplement with Detailed Searches (as needed), Step 4: Stack Guidelines

### Community 123 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 124 - "graphify reference: commit hook and native CLAUDE.md integration"
Cohesion: 0.50
Nodes (3): For git commit hook, For native CLAUDE.md integration, graphify reference: commit hook and native CLAUDE.md integration

### Community 125 - "graphify reference: incremental update and cluster-only"
Cohesion: 0.50
Nodes (3): For --cluster-only, For --update (incremental re-extraction), graphify reference: incremental update and cluster-only

### Community 144 - "md"
Cohesion: 0.67
Nodes (4): $type, $value, md, md

## Knowledge Gaps
- **992 isolated node(s):** `fs`, `path`, `fs`, `path`, `fs` (+987 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `primitive` connect `primitive` to `gray`, `color`, `spacing`, `radius`, `fontSize`, `shadow`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `color` connect `color` to `button`?**
  _High betweenness centrality (0.006) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `TailwindConfigGenerator` (e.g. with `TestGeneratedConfigIsValidJs` and `TestTailwindConfigGenerator`) actually correct?**
  _`TailwindConfigGenerator` has 2 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `DesignSystemGenerator` (e.g. with `TestDomainDetection` and `TestPersistence`) actually correct?**
  _`DesignSystemGenerator` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `fs`, `path`, `fs` to the rest of the system?**
  _992 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `gray` be split into smaller, more focused modules?**
  _Cohesion score 0.05370101596516691 - nodes in this community are weakly interconnected._
- **Should `BM25` be split into smaller, more focused modules?**
  _Cohesion score 0.06693877551020408 - nodes in this community are weakly interconnected._