module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: ["eslint:recommended"],
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
  },
  plugins: ["vue"],
  overrides: [
    {
      files: ["*.vue"],
      parser: "vue-eslint-parser",
    },
  ],
  rules: {
    // Add custom rules here
  },
};
