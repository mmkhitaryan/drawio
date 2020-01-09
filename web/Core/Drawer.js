import Connect from "./Connect.js";
import Point from "./Point.js";
import QueryDrawer from "./QueryDrawer.js";

export default class Drawer {
    /**
     * Drawer constructor
     * @param {Number} speed
     */
    constructor(speed) {
        this._speed = speed;
        this._isDrawing = false;
        this.x = 0;
        this.y = 0;
        this.points = [];
        this._canvas = document.querySelector("#can");
        this._qDrawer = new QueryDrawer(this._canvas);
        this.connect = new Connect(location.host, (data) => {
            this._qDrawer.addQuery(data);
        });
        this.initialize();
    }

    /**
     * Draw user line
     *
     * @param {Object} line
     * @returns {Uint16Array}
     */
    drawPoint(x, y) {
        const point = Point.encodePoint({x: x, y: y});
        this._qDrawer.addToLine(point);

        return point;
    }

    /**
     * Canvas initialize
     */
    initialize() {
        this._canvas.width  = window.innerWidth - 20;
        this._canvas.height = window.innerHeight - 20;

        this._canvas.addEventListener('mousedown', (event) => {
            this.x = event.pageX - this._canvas.offsetLeft;
            this.y = event.pageY - this._canvas.offsetTop;
            this.drawPoint(this.x, this.y);
            this._isDrawing = true;
        });

        this._canvas.addEventListener('mousemove', (event) => {
            if (this._isDrawing) {
                this.x = event.pageX - this._canvas.offsetLeft;
                this.y = event.pageY - this._canvas.offsetTop;
                let endpoint = drawPoint(this.x, this.y);
                this.points.push(endpoint);
            }
        }, false);

        this._canvas.addEventListener('mouseup', () => {
            this.x = 0;
            this.y = 0;
            this._isDrawing = false;
            this.connect.sendLine(this.points);
            this.points = [];
            this._qDrawer.clearLine();
        })
    }
}