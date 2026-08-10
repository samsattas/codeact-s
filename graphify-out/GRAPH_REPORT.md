# Graph Report - .  (2026-08-10)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 999 nodes · 1568 edges · 81 communities (62 shown, 19 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 82 edges (avg confidence: 0.68)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `585796ab`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Community 0
- Community 1
- Community 2
- Community 3
- Community 4
- Community 5
- Community 6
- Community 7
- Community 8
- Community 9
- Community 10
- Community 11
- Community 12
- Community 13
- Community 14
- Community 15
- Community 16
- Community 17
- Community 18
- Community 19
- Community 20
- Community 21
- Community 22
- Community 23
- Community 24
- Community 25
- Community 26
- Community 27
- Community 28
- Community 29
- Community 30
- Community 31
- Community 32
- Community 33
- Community 34
- Community 35
- Community 36
- Community 37
- Community 38
- Community 39
- Community 40
- Community 41
- Community 42
- Community 43
- Community 44
- Community 45
- Community 46
- Community 47
- Community 48
- Community 49
- Community 50
- Community 51
- Community 52
- Community 53
- Community 54
- Community 55
- Community 56
- Community 57
- Community 58
- Community 59
- Community 60
- Community 61
- Community 62
- Community 63
- Community 64
- Community 65
- Community 66
- Community 67
- Community 68
- Community 69
- Community 70
- Community 71
- Community 72
- Community 73
- Community 74
- Community 75
- Community 76
- Community 77
- Community 78
- Community 79
- Community 80

## God Nodes (most connected - your core abstractions)
1. `TailwindConfigGenerator` - 58 edges
2. `TestTailwindConfigGenerator` - 35 edges
3. `ShadcnInstaller` - 34 edges
4. `DesignSystemGenerator` - 29 edges
5. `TestShadcnInstaller` - 26 edges
6. `color` - 15 edges
7. `search()` - 15 edges
8. `search_with_context()` - 12 edges
9. `gray` - 12 edges
10. `spacing` - 12 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `FromEnv()`  [INFERRED]
  cmd/web/main.go → internal/llm/llm.go
- `main()` --calls--> `SetWorkDir()`  [INFERRED]
  cmd/web/main.go → internal/tools/tools.go
- `newTestServer()` --calls--> `SetWorkDir()`  [INFERRED]
  cmd/web/main_test.go → internal/tools/tools.go
- `main()` --calls--> `FromEnv()`  [INFERRED]
  cmd/agent/main.go → internal/llm/llm.go
- `main()` --calls--> `SetWorkDir()`  [INFERRED]
  cmd/agent/main.go → internal/tools/tools.go

## Import Cycles
- None detected.

## Communities (81 total, 19 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.05
Nodes (53): $type, $value, $type, $value, $type, $value, $type, $value (+45 more)

### Community 1 - "Community 1"
Cohesion: 0.07
Nodes (42): BM25, detect_domain(), get_cip_brief(), _load_csv(), Load CSV and return list of dicts, Core search function using BM25, Auto-detect the most relevant domain from query, Main search function with auto-domain detection (+34 more)

### Community 2 - "Community 2"
Cohesion: 0.09
Nodes (38): Agent, fakeProvider, Outcome, Step, fatal(), main(), repl(), runTask() (+30 more)

### Community 3 - "Community 3"
Cohesion: 0.06
Nodes (30): BM25, detect_domain(), _domain_keywords(), _get_bm25(), _load_csv(), _load_product_keywords(), _normalize(), Apply synonym substitution before tokenizing. (+22 more)

### Community 4 - "Community 4"
Cohesion: 0.09
Nodes (36): format_context(), format_result(), main(), Format a single search result for display, Format contextual recommendations for display., BM25, calculate_pattern_break(), detect_domain() (+28 more)

### Community 5 - "Community 5"
Cohesion: 0.05
Nodes (37): $type, $value, background, destructive, destructive-foreground, foreground, muted, muted-foreground (+29 more)

### Community 6 - "Community 6"
Cohesion: 0.06
Nodes (34): $type, $value, $type, $value, $type, $value, $type, $value (+26 more)

### Community 7 - "Community 7"
Cohesion: 0.14
Nodes (24): get_context(), is_allowed_exception(), is_allowed_rgba(), is_inside_block(), load_css_variables(), main(), print_result(), print_summary() (+16 more)

### Community 8 - "Community 8"
Cohesion: 0.07
Nodes (14): Test adding colors multiple times., Test adding full color palette., Test TailwindConfigGenerator class., Test initialization with default settings., Test generating config with custom colors., Test validating valid configuration., Test validating config with empty theme extensions., Test writing configuration to file. (+6 more)

### Community 9 - "Community 9"
Cohesion: 0.12
Nodes (19): BM25, detect_domain(), _load_csv(), Load CSV and return list of dicts, Core search function using BM25, Auto-detect the most relevant domain from query, Main search function with auto-domain detection, Search across all domains and combine results (+11 more)

### Community 10 - "Community 10"
Cohesion: 0.26
Nodes (20): CountLinesByExt(), FileTree(), Grep(), ListDir(), ReadFile(), resolve(), RunCommand(), skipDir() (+12 more)

### Community 11 - "Community 11"
Cohesion: 0.16
Nodes (15): main(), Context, T, newTestServer(), TestHandleInfo(), TestHandleRunRejectsMissingTask(), TestHandleRunRejectsNonPost(), TestHandleRunStreamsFailureAfterAttempts() (+7 more)

### Community 12 - "Community 12"
Cohesion: 0.15
Nodes (19): _e(), generate_chart_slide(), generate_cta_slide(), generate_deck(), generate_metrics_slide(), generate_problem_slide(), generate_solution_slide(), generate_testimonial_slide() (+11 more)

### Community 13 - "Community 13"
Cohesion: 0.17
Nodes (8): DesignSystemGenerator, Generates design system recommendations from aggregated searches., Load reasoning rules from CSV., Find matching reasoning rule for a category., Apply reasoning rules to search results., TestReasoningMatch, The exact reproduction from issue #428., TestEndToEndCoherence

### Community 14 - "Community 14"
Cohesion: 0.11
Nodes (10): Generate Tailwind CSS configuration files., Add full color palette (50-950 shades) for a base color. Args: name: Color name…, TailwindConfigGenerator, Test adding custom spacing., Test generating TypeScript configuration., Test generating JavaScript configuration., Test generating config with plugins., Test initialization for JavaScript config. (+2 more)

### Community 15 - "Community 15"
Cohesion: 0.17
Nodes (17): generate_css_for_background(), get_background_image(), get_curated_images(), get_overlay_css(), get_pexels_search_url(), load_backgrounds_config(), load_brand_colors(), main() (+9 more)

### Community 16 - "Community 16"
Cohesion: 0.20
Nodes (15): apply_color(), apply_viewbox_size(), extract_svgs(), generate_batch(), generate_icon(), generate_sizes(), load_env(), main() (+7 more)

### Community 17 - "Community 17"
Cohesion: 0.12
Nodes (16): $type, $value, $type, $value, $type, $value, $type, $value (+8 more)

### Community 18 - "Community 18"
Cohesion: 0.16
Nodes (7): Test adding components in dry run mode., Test ShadcnInstaller class., Test initialization with dry run mode., Test checking for existing shadcn config., Test getting installed components when none exist., Test getting installed components without config., TestShadcnInstaller

### Community 19 - "Community 19"
Cohesion: 0.13
Nodes (8): main(), Add custom font families. Args: fonts: Dict of font_type: [font_names] e.g.,…, Add custom spacing values. Args: spacing: Dict of name: value e.g., {'18':…, Add custom breakpoints. Args: breakpoints: Dict of name: width e.g., {'3xl':…, Add plugin requirements. Args: plugins: List of plugin names e.g.,…, Get plugin recommendations based on configuration. Returns: List of recommended…, Validate configuration. Returns: Tuple of (valid, message), Add custom colors to theme. Args: colors: Dict of color_name: color_value Value…

### Community 20 - "Community 20"
Cohesion: 0.27
Nodes (7): _query_wants_dark(), True when a styles.csv row describes itself as dark-first., True when the query explicitly asks for a dark theme., Resolve the mode the rest of the output has to agree with., _resolve_color_mode(), _style_is_dark_primary(), TestModeResolution

### Community 21 - "Community 21"
Cohesion: 0.27
Nodes (5): _palette_is_dark(), WCAG relative luminance of a #RRGGBB string, or None if unparseable., True when a colors.csv row's Background is a dark surface., _relative_luminance(), TestLuminance

### Community 22 - "Community 22"
Cohesion: 0.22
Nodes (11): calculateCompliance(), colorDistance(), displayPalette(), extractHexColors(), findNearestBrandColor(), fs, generateImageMagickCommand(), hexToRgb() (+3 more)

### Community 23 - "Community 23"
Cohesion: 0.25
Nodes (13): checkManifest(), formatBytes(), formatOutput(), fs, main(), parseFilename(), path, RULES (+5 more)

### Community 24 - "Community 24"
Cohesion: 0.18
Nodes (11): format_markdown(), format_master_md(), generate_design_system(), persist_design_system(), Format design system as markdown., Main entry point for design system generation. Args: query: Search query (e.g.,…, Slugify a name into a single safe path segment. Only [a-z0-9_-] survives; every…, Persist design system to design-system/<project>/ folder using Master +… (+3 more)

### Community 25 - "Community 25"
Cohesion: 0.18
Nodes (13): $type, $value, border, padding, radius, shadow, border, card (+5 more)

### Community 26 - "Community 26"
Cohesion: 0.20
Nodes (7): main(), Add all available shadcn/ui components. Args: overwrite: If True, overwrite…, List installed components. Returns: Tuple of (success, message with component…, Check if shadcn is initialized in project. Returns: True if components.json…, Get list of already installed components. Returns: List of installed component…, Read shadcn version from project package.json; fall back to a pinned default., Add shadcn/ui components. Args: components: List of component names to add…

### Community 27 - "Community 27"
Cohesion: 0.15
Nodes (8): Path, Handle shadcn/ui component installation., Initialize installer. Args: project_root: Project root directory (default:…, ShadcnInstaller, Test adding components without shadcn config., Test listing installed components when none exist., Test initialization with custom project root., Test checking for non-existent shadcn config.

### Community 28 - "Community 28"
Cohesion: 0.31
Nodes (11): diagnosisHint(), fenced(), retryPrompt(), T, TestDiagnosisHintExitStatus(), TestDiagnosisHintUndefined(), TestDiagnosisHintUnknownErrorIsEmpty(), TestDiagnosisHintUnusedImport() (+3 more)

### Community 29 - "Community 29"
Cohesion: 0.24
Nodes (11): extensions, formatReport(), fs, getFiles(), main(), parseArgs(), path, patterns (+3 more)

### Community 30 - "Community 30"
Cohesion: 0.20
Nodes (8): Tests for tailwind_config_gen.py, Reduce a generated TS/JS config to a bare assignable object so it can be handed…, Regression guard for the missing-comma bug between the ``theme`` block and…, The property preceding ``plugins`` must end with a comma (pure-Python check, so…, The emitted config parses as valid JS via ``node --check``., _strip_to_object(), TestGeneratedConfigIsValidJs, parametrize

### Community 31 - "Community 31"
Cohesion: 0.20
Nodes (6): Generate configuration file content. Returns: Configuration file as string, Generate TypeScript configuration., Generate JavaScript configuration., Format plugins array for config. Validates each plugin name against a strict…, Add indentation to JSON string., Write configuration to file. Returns: Tuple of (success, message)

### Community 32 - "Community 32"
Cohesion: 0.31
Nodes (10): extractColorsFromTable(), extractCoreAttributes(), extractHexColors(), extractImageStyle(), extractTypography(), extractVoice(), fs, generatePromptAddition() (+2 more)

### Community 33 - "Community 33"
Cohesion: 0.18
Nodes (8): args, fs, minimal, MINIMAL_TOKENS, path, projectRoot, tokensPath, wrapStyle

### Community 34 - "Community 34"
Cohesion: 0.18
Nodes (11): fast, normal, slow, $type, $value, $type, $value, primitive (+3 more)

### Community 35 - "Community 35"
Cohesion: 0.18
Nodes (6): Test adding components with overwrite flag., Test successful component addition., Test component addition with subprocess error., Test component addition when npx is not found., Test successful addition of all components., patch

### Community 36 - "Community 36"
Cohesion: 0.25
Nodes (7): Client, Context, NewOllamaProvider(), ollamaChatRequest, ollamaChatResponse, ollamaMessage, OllamaProvider

### Community 37 - "Community 37"
Cohesion: 0.22
Nodes (6): Any, Path, Initialize generator. Args: typescript: If True, generate .ts config, else .js…, Determine default output path., Create base configuration structure., Get default content paths for framework.

### Community 38 - "Community 38"
Cohesion: 0.29
Nodes (9): enhance_prompt(), generate_batch(), generate_logo(), load_env(), main(), Enhance the logo prompt with style and industry modifiers, Generate a logo using Gemini models with image generation Args: aspect_ratio:…, Generate multiple logo variants with different styles (+1 more)

### Community 39 - "Community 39"
Cohesion: 0.36
Nodes (9): flattenTokens(), fs, generateCSS(), generateTailwind(), main(), parseArgs(), path, resolveReference() (+1 more)

### Community 40 - "Community 40"
Cohesion: 0.20
Nodes (10): fg, font-size, hover-bg, button, $type, $value, $type, $value (+2 more)

### Community 41 - "Community 41"
Cohesion: 0.24
Nodes (7): Client, Context, NewAnthropicProvider(), anthropicMessage, AnthropicProvider, anthropicRequest, anthropicResponse

### Community 42 - "Community 42"
Cohesion: 0.33
Nodes (8): adjustBrightness(), { execFileSync }, extractColorsFromMarkdown(), fs, generateColorScale(), main(), path, updateDesignTokens()

### Community 43 - "Community 43"
Cohesion: 0.28
Nodes (8): Path, Regression tests for validate-tokens.cjs. The validator used to skip any line…, A hardcoded hex on the same line as a var() token is still a violation., A line that references only tokens produces no false positives., _run(), test_flags_hardcoded_hex_sharing_line_with_token(), test_token_only_line_reports_no_violation(), CompletedProcess

### Community 44 - "Community 44"
Cohesion: 0.29
Nodes (8): padding-x, input, $type, $value, focus-ring, padding-x, $type, $value

### Community 45 - "Community 45"
Cohesion: 0.25
Nodes (8): $type, $value, $type, $value, semantic, spacing, component, section

### Community 46 - "Community 46"
Cohesion: 0.29
Nodes (8): $type, $value, $type, $value, radius, default, full, default

### Community 47 - "Community 47"
Cohesion: 0.25
Nodes (8): ansi_ljust(), format_ascii_box(), hex_to_ansi(), Convert hex color to ANSI True Color swatch (██) with fallback., Like str.ljust but accounts for zero-width ANSI escape sequences., Create a Unicode section separator: ├─── NAME ───...┤, Format design system as Unicode box with ANSI color swatches., section_header()

### Community 48 - "Community 48"
Cohesion: 0.43
Nodes (3): _filter_anti_patterns_for_mode(), Drop "avoid dark mode" advice once dark mode is the resolved answer., TestAntiPatternGating

### Community 49 - "Community 49"
Cohesion: 0.43
Nodes (3): Pick the highest-ranked palette matching the resolved mode. Only the dark case…, _select_palette_for_mode(), TestPaletteSelection

### Community 50 - "Community 50"
Cohesion: 0.47
Nodes (6): sm, shadow, sm, sm, $type, $value

### Community 51 - "Community 51"
Cohesion: 0.33
Nodes (6): _detect_page_type(), format_page_override_md(), _generate_intelligent_overrides(), Format a page-specific override file with intelligent AI-generated content., Generate intelligent overrides based on page type using layered search. Uses…, Detect page type from context and search results.

### Community 52 - "Community 52"
Cohesion: 0.40
Nodes (4): component, dark, semantic, $schema

### Community 53 - "Community 53"
Cohesion: 0.60
Nodes (5): $type, $value, bg, bg, bg

### Community 54 - "Community 54"
Cohesion: 0.60
Nodes (5): lg, $type, $value, lg, lg

### Community 55 - "Community 55"
Cohesion: 0.67
Nodes (4): padding-y, padding-y, $type, $value

### Community 56 - "Community 56"
Cohesion: 0.67
Nodes (4): radius, radius, $type, $value

### Community 57 - "Community 57"
Cohesion: 0.67
Nodes (4): xl, xl, $type, $value

### Community 58 - "Community 58"
Cohesion: 0.67
Nodes (4): $type, $value, md, md

### Community 59 - "Community 59"
Cohesion: 0.67
Nodes (4): $type, $value, none, none

### Community 60 - "Community 60"
Cohesion: 0.83
Nodes (3): _check_file(), main(), _read_rows()

### Community 64 - "Community 64"
Cohesion: 0.20
Nodes (6): Execute searches across multiple domains., Select best matching result based on priority keywords., Extract results list from search result dict., Generate complete design system recommendation. variance/motion/density are…, Bucket a 1-10 dial value into its tier config. Returns None if value is None., _resolve_dial()

## Knowledge Gaps
- **123 isolated node(s):** `fs`, `path`, `fs`, `path`, `fs` (+118 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `primitive` connect `Community 34` to `Community 0`, `Community 6`, `Community 46`, `Community 17`, `Community 50`, `Community 52`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **Why does `color` connect `Community 0` to `Community 34`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `semantic` connect `Community 45` to `Community 52`, `Community 5`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `TailwindConfigGenerator` (e.g. with `TestGeneratedConfigIsValidJs` and `TestTailwindConfigGenerator`) actually correct?**
  _`TailwindConfigGenerator` has 2 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `DesignSystemGenerator` (e.g. with `TestDomainDetection` and `TestPersistence`) actually correct?**
  _`DesignSystemGenerator` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `fs`, `path`, `fs` to the rest of the system?**
  _123 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.05370101596516691 - nodes in this community are weakly interconnected._