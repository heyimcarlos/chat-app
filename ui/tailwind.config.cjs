module.exports = {
    content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
    darkMode: 'class',
    theme: {
        extend: {
            screens: {
                xs: '450px'
            },
            colors: {
                'custom-teal': '#5de4c7',
                'custom-zinc': '#202020',
                'custom-black': '#101010',
                'custom-blue': '#0069C2',
                'custom-white': '#f3f4f6'
            },
            fontFamily: {
                // @info: Adding a utility class for the font
                mplus: ["'M PLUS Rounded 1c'", 'Verdana', 'sans-serif']
            },
            keyframes: {
                wiggle: {
                    '0%, 100%': { transform: 'rotate(-3deg)' },
                    '50%': { transform: 'rotate(3deg)' }
                },
                shake: {
                    '0%': { transform: 'translate(0rem)' },
                    '25%': { transform: 'translateX(5px)' },
                    '75%': { transform: 'translateX(-5px)' },
                    '100%': { transform: 'translate(0rem)' }
                }
            },
            animation: {
                wiggle: 'wiggle 0.2s linear 0s 1 normal none running',
                shake: 'shake 0.2s linear 0s 1 normal none running'
            }
        }
    },
    plugins: [require('@tailwindcss/forms')],
};
