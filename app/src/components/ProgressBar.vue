<template>
    <div class="progress-bar">
        <div class="bar-fill" :style="{'height': dynamicHeight}"></div>
        <div class="bar-status">
            <v-row no-gutters >
                <v-col class="bar-status-label">
                    {{ blocked ? 'Restoring' : dynamicHeight }}
                </v-col>
                <v-col v-if="blocked" class="bar-status-gif">
                    <v-img height="20px" width="20px" :src="loadingImage"></v-img>
                </v-col>
            </v-row></div>
    </div>
</template>

<script>
export default {
    props: {
        blocked: Boolean,
        progress: Number
    },
    data: () => ({
        max: 240,
        loadingImage: require('@/assets/refresh.gif')
    }),
    computed: {
        dynamicHeight() {
            let height = (100 - parseInt(((this.progress / this.max) * 100), 10)) + '%';
            //eslint-disable-next-line
            console.log(height);
            return height;
        }
    }
}
</script>

<style>
.progress-bar {
    position: fixed;
    top: 5px;
    left: 5px;
    height: 240px;
    width: 20px;
    background-color: aliceblue;
    border: 2px solid #222;
}
.bar-fill {
    position: absolute;
    bottom: 1px;
    left: -2px;
    background-color: #222;
    width: 20px;
    transition: 200ms linear;
}
.bar-status {
    position: absolute;
    width: 150px;
    top: 0;
    left: 25px;
    font-size: 14pt;
    color: #222;
    background-color: aliceblue;
}
.bar-status-label {
    padding: 4px !important;
}
.bar-status-gif {
    padding: 8px !important;
}
</style>