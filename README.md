# CHAOSS AI Detection Tool

This project is intended to aid the open source community by detecting common ways that agents identify themselves for the purposes of aiding already-overworked maintainers and providing a tool for projects that track community health of open source repositories.

This project will likely eventually consist of two parts:

1. A golang project that provides the core detection capabilities, suitable for potential use by other projects
2. A github action making use of this detection core and provides easy, sensible defaults to allow maintainers to:
   - Automate certain basic tasks related to enforcing AI policy, such as applying labels to disclosed/detected AI content or closing pull requests
   - Enable AI usage in their projects to be more easily displayed by community health and metrics tools (such as the [Augur](https://github.com/chaoss/augur/) or [GrimoireLab](https://github.com/chaoss/grimoirelab/) projects)


## Project Philosophy

TODO


