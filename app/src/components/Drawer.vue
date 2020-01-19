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
        drawer = new Drawer(2, location.host /*'157.245.120.30'*/, () => {
            this.$emit('draw');
        });
    },
    methods: {
        blockBrush() {
            //eslint-disable-next-line
            console.log('block');
            drawer.block();
        },
        unblockBrush() {
            //eslint-disable-next-line
            console.log('unblock');
            drawer.unblock();
        },
        setColor(color) {
            drawer.setColor(color);
        }
    }
}
</script>