# Asset Pipeline

Visual engines in v1:

- Mermaid for topology, sequence, ER, state, class, and C4-like diagrams.
- SVG for precise custom diagrams, annotations, and connectors.
- Charts for quantitative evidence.
- Image generation for raster illustrations, editorial imagery, mockups,
  textures, and concept visuals when explicitly justified by the report plan.

Choose the engine by job. Do not use image generation for visuals that should
be deterministic SVG, HTML/CSS, Mermaid, or chart code.

If an imagegen asset is needed, decide at the start of the report and dispatch
it in parallel because image generation is slow in Codex. The image worker owns
only the raster asset and returns path, prompt, metadata, and fit notes. It must
not edit report HTML or profile files.

Record assets compactly in `canvas.json.assets`.
