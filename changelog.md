<p align="center">
  <picture>
    <source srcset="https://raw.githubusercontent.com/svengreb/assets/main/static/images/logos/projects/home-ops/heroes/dark.svg" width="100%" media="(prefers-color-scheme: light), (prefers-color-scheme: no-preference)" />
    <source srcset="https://raw.githubusercontent.com/svengreb/assets/main/static/images/logos/projects/home-ops/heroes/light.svg" width="100%" media="(prefers-color-scheme: dark)" />
    <img src="https://raw.githubusercontent.com/svengreb/assets/main/static/images/logos/projects/home-ops/heroes/dark.svg" width="100%" />
  </picture>
</p>

<p align="center">
  <picture>
    <source srcset="https://raw.githubusercontent.com/svengreb/assets/main/static/images/elements/separators/logo/footer/dark/spaced.svg?sanitize=true" width="100%" media="(prefers-color-scheme: light), (prefers-color-scheme: no-preference)" />
    <source srcset="https://raw.githubusercontent.com/svengreb/assets/main/static/images/elements/separators/logo/footer/light/spaced.svg?sanitize=true" width="100%" media="(prefers-color-scheme: dark)" />
    <img src="https://raw.githubusercontent.com/svengreb/assets/main/static/images/elements/separators/logo/footer/dark/spaced.svg?sanitize=true" width="100%" />
  </picture>
</p>

<!--lint disable no-duplicate-headings-->

# 0.1.0

![Release Date: 2025-11-23](https://img.shields.io/static/v1.svg?style=flat-square&label=Release%20Date&message=2025-11-23&colorA=4c566a&colorB=88c0d0) [![Project Board](https://img.shields.io/static/v1.svg?style=flat-square&label=Project%20Board&message=0.1.0&logo=github&logoColor=eceff4&colorA=4c566a&colorB=88c0d0)](https://github.com/users/svengreb/projects/14/views/10) [![Milestone](https://img.shields.io/static/v1.svg?style=flat-square&label=Milestone&message=0.1.0&logo=github&logoColor=eceff4&colorA=4c566a&colorB=88c0d0)](https://github.com/svengreb/home-ops-cli/milestone/1)

⇅ [Show all commits][111]

## Features

<details>
<summary><strong>Initial repository setup</strong> — #1 (⊶ [4bd74e93][6]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Set up the basic project tooling and structure which includes…

- formatting and linting tools like [Prettier][1], [lint-staged][2] and [remark][3].
- essential _Git_ configuration files.
- [EditorConfig][4] configuration for cross-IDE settings.

</details>

<details>
<summary><strong>Support package for extracting archive files</strong> — #2 (⊶ [7d8194b0][5]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Implemented a support package to simplify extracting archive files (e.g., `.zip`, `.tar`, `.tar.gz`, `.tar.xz`, etc.) using the [github.com/mholt/archives][7] module. This helps to reduce duplication across the CLI codebase when using the module.
To ensure safe extraction a ["ZipSlip" vulnerability][8] [^1] (path traversal) protection and proper directory creation has been implemented.

</details>

<details>
<summary><strong>Basic HTTP REST client for the Home Assistant API</strong> — #3 (⊶ [20e6be5c][13]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Implemented a lightweight HTTP [REST][9] client to interact with the [Home Assistant HTTP API][10] [^2], providing a clean, reusable interface for issuing authenticated requests from the CLI.
To simplify the implementation the [github.com/go-resty/resty][11] [^3] (v2) module is used and for human-facing logs the [github.com/charmbracelet/log][12] module to "prettify" the output.

</details>

<details>
<summary><strong>Extended Home Assistant API REST client to support entity states and actions</strong> — #4 (⊶ [cfdd62d3][14]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Extended the Home Assistant REST client implemented in #3 to retrieve [entity states][15] [^4] [^5] and trigger entity ["Actions"][16] [^6] [^7] (previously known as "Services"), enabling foundational integration with Home Assistant devices and automations.

</details>

<details>
<summary><strong>Extended Home Assistant API REST client to support types for "Backup" integration</strong> — #5 (⊶ [73b23dbc][17]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Extended the Home Assistant REST client implemented in #3 to support types for the (built-in) ["Backup" integration][18], enabling programmatic parsing and processing of [the `backup.json` file][19] that is included in the archive files created by the integration. This allows for processing the archive files in custom backup automations, e.g. to extract only specific data.

</details>

<details>
<summary><strong>Initial implementation for <code>homeopsctl</code> main CLI application</strong> — #6 (⊶ [bf8c977d][20]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Created the foundational structure for the root command of the `homeopsctl` application using the [github.com/spf13/cobra][21] module, providing a starting point for (interactive) commands like interactions with Home Assistant via its API through the custom REST client (implemented in #3).
Initially the CLI accepts global flags to toggle debug logging, change the output format and a some other basic flags. All configurations can be parsed as environment variables using the `HOMEOPSCTL_` prefix.

</details>

<details>
<summary><strong>Subcommand to extract specific data from automated Home Assistant backup archives</strong> — #7 (⊶ [0cafe01b][22]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Implemented a [subcommand][23] that creates custom, lightweight backups of Home Assistant by extracting only selected data from existing backup archives that are created automatically by [the built-in "backup" integration][24].
When used within the scope of [the HomeOps GitOps & IaC repository][25] for [Kubernetes][26], the extracted data will be written into a dedicated [Kubernetes PVC][5] so that backup solutions like [K8up][27] can pick it up for scheduled backups.

</details>

<details>
<summary><strong>Build and release process with "ko" and automated container releases to GHCR via GitHub Actions</strong> — #8 (⊶ [ad6b6a4b][29]) by <a href="https://github.com/svengreb" target="_blank" rel="noreferrer">@svengreb</a></summary>

↠ Set up [containerized][30] builds of the CLI using [ko][31], and automate image publishing to the [GitHub Container Registry (GHCR)][35] through [GitHub Actions][36]. This enables consistent, reproducible images and streamlined releases.
The [configuration for ko][32] is persisted through [a `.ko.yaml` file][33] in the root of the repository which will be picked up during the build process. The [`ko-build/setup-ko` GitHub Action][34] is used to run the whole process in an automated fashion for development states, requested explicitly though repository maintainers, and tagged releases.

</details>

<p align="center">Copyright &copy; 2018 <a href="https://www.svengreb.de" target="_blank">Sven Greb</a></p>

<p align="center">
  <a href="https://github.com/svengreb/home-ops-cli/blob/main/license" target="_blank">
    <img src="https://img.shields.io/static/v1.svg?style=flat-square&label=License&message=Apache%202.0&logoColor=eceff4&logo=creativecommons&colorA=4c566a&colorB=88c0d0"/>
  </a>
    <a href="https://www.svengreb.de">
    <img src="https://img.shields.io/static/v1.svg?style=flat-square&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAABMklEQVQ4jcWQvUoDQRRGz52s5IfVIiDWPkGKFFaCIVaGdIagjcFAwICFb7DvIK6QQlNpY2UQLMQVBbEQ0SewFkGbKCQmOzaTJay7/lR+zTAf9xwuF/47Mv45rdezqWEq72v/RWZnHgqOMwDwHMfSj085JSqb6Pu38we7r18E3nqzhmYbsE11rxKsAvhDfQiSM30XYbOw57YDwfnaRl6U3ABWaMNn806H+oGPzBX3d+4UgChZiYBHYBgGsBLoKoAyhR0x9G20Zmpc4P1ZoMQDcwMNclFrdhBKv6M5WWi7ZQGtjEUn35IV4OwnVjSX/WGmKqCDDUa5rmyle3bvGFiMg3WGUsF1u0EXHoqTRMGRgkAy2eugKZrqijRLYThWANBpNDL2h3UE0J0YLJdbrfe42f/NJ0wqY7/KcXKPAAAAAElFTkSuQmCC&label=lovely%20crafted%20in&message=Germany&colorA=4c566a&colorB=88c0d0" />
  </a>
</p>

<!--
+------------------+
+ Formatting Notes +
+------------------+

The `<summary />` tag must be separated with a blank line from the actual item content paragraph,
otherwise Markdown elements are not parsed and rendered!

+------------------+
+ Symbol Reference +
+------------------+
↠ (U+21A0): Start of a log section description
— (U+2014): Separator between a log section title and the metadata
⇄ (U+21C4): Separator between a issue ID and pull request ID in a log metadata
⊶ (U+22B6): Icon prefix for the short commit SHA checksum in a log metadata
⇅ (U+21C5): Icon prefix for the link of the Git commit history comparison on GitHub
-->

<!--lint disable final-definition-->

<!-- Base -->

<!-- Shared -->

<!-- 0.1.0 -->

[1]: https://prettier.io
[2]: https://github.com/lint-staged/lint-staged
[3]: https://remark.js.org
[4]: https://editorconfig.org
[5]: https://github.com/svengreb/home-ops-cli/commit/7d8194b0eebf1ddd45f363aeb241edf1d7f8864a
[6]: https://github.com/svengreb/home-ops-cli/commit/4bd74e93d09c7015a264ee195d8563ec0cc42003
[7]: https://github.com/mholt/archives
[8]: https://security.snyk.io/research/zip-slip-vulnerability
[9]: https://en.wikipedia.org/wiki/REST
[10]: https://developers.home-assistant.io/docs/api/rest
[11]: https://github.com/go-resty/resty
[12]: https://github.com/charmbracelet/log
[13]: https://github.com/svengreb/home-ops-cli/commit/20e6be5c75b74f01d142727139f2696cc72e0f8a
[14]: https://github.com/svengreb/home-ops-cli/commit/cfdd62d39d1cb3f9a70ff7c77ff653d4428803f0
[15]: https://www.home-assistant.io/docs/configuration/state_object/#about-the-state-object
[16]: https://www.home-assistant.io/docs/automation/action
[17]: https://github.com/svengreb/home-ops-cli/commit/73b23dbc44969631e5de0f4545e1a336c1353b29
[18]: https://www.home-assistant.io/integrations/backup
[19]: https://github.com/home-assistant/core/blob/e572f8d4/homeassistant/components/backup/util.py#L74-L76
[20]: https://github.com/svengreb/home-ops-cli/commit/bf8c977d393969dca0c7fa6bfc249881ba924683
[21]: https://github.com/spf13/cobra
[22]: https://github.com/svengreb/home-ops-cli/commit/0cafe01b631c01fb9b9c4161154b9c8973cf141f
[23]: https://cobra.dev/docs/how-to-guides/working-with-commands
[24]: https://www.home-assistant.io/integrations/backup
[25]: https://github.com/svengreb/home-ops
[26]: https://kubernetes.io
[27]: https://kubernetes.io/docs/concepts/storage/persistent-volumes
[28]: https://k8up.io
[29]: https://github.com/svengreb/home-ops-cli/commit/ad6b6a4b1f9ed30855d59098904b03257227b91e
[30]: https://opencontainers.org
[31]: https://ko.build
[32]: https://ko.build/configuration
[33]: https://github.com/svengreb/home-ops-cli/blob/ad6b6a4b/.ko.yaml
[34]: https://github.com/ko-build/setup-ko
[35]: https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry
[36]: https://github.com/features/actions

<!-- prettier-ignore-start -->

[^1]: https://github.com/snyk/zip-slip-vulnerability
[^2]: https://www.home-assistant.io/integrations/api
[^3]: https://resty.dev
[^4]: https://data.home-assistant.io/docs/states
[^5]: https://developers.home-assistant.io/docs/dev_101_states
[^6]: https://www.home-assistant.io/docs/scripts/perform-actions
[^7]: https://data.home-assistant.io/docs/services

<!-- prettier-ignore-end -->
