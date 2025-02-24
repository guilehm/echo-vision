import { dirname } from "path";
import { fileURLToPath } from "url";
import { FlatCompat } from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const eslintConfig = [
  ...compat.extends("next"),
  {
    rules: {
      "newline-per-chained-call": ["error", { ignoreChainWithDepth: 4 }],
      indent: ["error", 2, { SwitchCase: 1 }],
      "linebreak-style": ["error", "unix"],
      quotes: ["error", "double"],
      // semi: ["error", "never"],
    },
  },
];

export default eslintConfig;
