<template>
        <canvas id="can" width="700" height="700"></canvas>
</template>

<script>
import Drawer from '@/drawer/Core/Drawer.js';

//eslint-disable-next-line
let drawer = null;

export default {
    data: () => ({
        speed: 2
    }),
    mounted() {
        drawer = new Drawer(2, '157.245.120.30', () => {
            this.$emit('draw');
        });
        window.scrollTo(2500 - ( window.innerWidth / 2 ), 2500 - ( window.innerHeight / 2 ));
    },
    methods: {
        blockBrush() {
            drawer.block();
        },
        unblockBrush() {
            drawer.unblock();
        },
        setColor(color) {
            drawer.setColor(color);
        }
    }
}
</script>