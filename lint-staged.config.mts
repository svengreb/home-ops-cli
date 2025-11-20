/*
 * Copyright 2018 Sven Greb <development@svengreb.de>
 * This source code is licensed under the Apache License 2.0 found in the license file.
 */

/**
 * Configurations for lint-staged.
 * @see https://github.com/okonet/lint-staged#configuration
 */
const config = {
  "*.{js,json,yaml,yml}": "prettier --check",
  "*.md": ["remark --no-stdout", "prettier --check"],
};

export default config;
