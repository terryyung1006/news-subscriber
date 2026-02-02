# General Coding Principles

## Reuse Existing Code
Do not create a new file if there is existing code that can achieve the same functionality. Always check for existing utilities, functions, or scripts before adding new ones.

## Necessity of New Scripts
Only add new scripts if they are essential for the repository and intended to be committed to Git. Avoid creating files for one-off operations; use temporary methods or run commands directly if possible.

## Clean Repository Hierarchy
Ensure any new file is placed in the appropriate directory. Maintain a clean and logical project structure. Do not clutter the root or unrelated directories.

## Limit Documentation
Avoid creating excessive documentation files. Use the `spec/` structure for features, and keep standalone docs to a minimum.

## No Example/Template Files
Do not create example files, template files, or "how-to" documentation unless explicitly requested by the user. When answering questions about how things work, provide explanations in the conversation rather than creating new files in the codebase.
