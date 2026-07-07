import reactHooks from 'eslint-plugin-react-hooks'
import tseslint from 'typescript-eslint'

// typescript-eslint recommended-type-checked plus the react-hooks rules.
// Type-aware linting covers src/, which tsconfig.json owns; tests and the
// config files run the untyped recommended set.
export default tseslint.config(
  { ignores: ['dist/'] },
  {
    files: ['src/**/*.{ts,tsx}'],
    extends: [...tseslint.configs.recommendedTypeChecked, reactHooks.configs.flat.recommended],
    languageOptions: {
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
  {
    files: ['tests/**/*.{ts,tsx}'],
    extends: [...tseslint.configs.recommended, reactHooks.configs.flat.recommended],
  },
  {
    files: ['*.ts', '*.js'],
    extends: [...tseslint.configs.recommended],
  },
)
