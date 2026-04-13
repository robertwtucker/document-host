/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import nextCoreWebVitals from 'eslint-config-next/core-web-vitals'
import nextTypescript from 'eslint-config-next/typescript'
import prettierConfig from 'eslint-config-prettier'
import tailwindcss from 'eslint-plugin-tailwindcss'

const config = [
  ...nextCoreWebVitals,
  ...nextTypescript,
  ...tailwindcss.configs['flat/recommended'],
  prettierConfig,
  {
    settings: {
      tailwindcss: {
        callees: ['cn'],
        entryPoint: 'app/globals.css',
      },
      next: {
        rootDir: ['./'],
      },
    },
    rules: {
      '@next/next/no-html-link-for-pages': 'off',
      'tailwindcss/no-custom-classname': 'off',
      '@typescript-eslint/no-empty-object-type': 'off',
      '@typescript-eslint/no-unused-vars': 'warn',
      'tailwindcss/classnames-order': 'warn',
    },
  },
]

export default config
