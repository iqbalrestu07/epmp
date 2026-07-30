/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        orange: {
          DEFAULT: "#ff6711",
          bright: "#ff8200",
        },
        glass: {
          bg: "rgba(255, 255, 255, 0.06)",
          bgStrong: "rgba(255, 255, 255, 0.09)",
          border: "rgba(255, 255, 255, 0.10)",
          borderStrong: "rgba(255, 255, 255, 0.16)",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        border: "rgba(255, 255, 255, 0.1)",
        input: "rgba(255, 255, 255, 0.1)",
        ring: "#ff6711",
        muted: {
          DEFAULT: "rgba(242, 239, 233, 0.62)",
          foreground: "rgba(242, 239, 233, 0.62)",
        },
      },
      fontFamily: {
        display: ["DM Serif Display", "Georgia", "serif"],
        body: ["Inter", "sans-serif"],
      },
      boxShadow: {
        glow: "0 10px 36px rgba(255, 103, 17, 0.35), 0 0 0 1px rgba(255, 255, 255, 0.10) inset, 0 -1px 0 rgba(0, 0, 0, 0.18) inset",
        glowHover: "0 14px 44px rgba(255, 103, 17, 0.45), 0 0 0 1px rgba(255, 255, 255, 0.14) inset",
      },
    },
  },
  plugins: [],
};
