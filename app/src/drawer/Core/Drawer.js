import Connect from "./Connect.js";
import Point from "./Point.js";
import Line from "./Line.js";
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
        this.line = new Line();
        this._canvas = document.querySelector("#can");
        this._qDrawer = new QueryDrawer(this._canvas);
        this.connect = new Connect("127.0.0.1:80", (data) => {
            this._qDrawer.addQuery(data);
        });
        this.initialize();
    }

    /**
     * Set line color
     * @param {String} color
     */
    setColor(color) {
        this.line.setColor(color);
        this.cut();

        return this;
    }

    /**
     * Cut line on 2 pieces
     */
    cut() {
        if(this._isDrawing) {
            const point = new Point(this.x, this.y);
            this.connect.sendLine(this.line.serialize());
            this.line = new Line([], this.line.getColor());
            this._qDrawer.clearLine(this.line.getColor());
            this.drawPoint(point);
        }

        return this;
    }

    /**
     * Draw user point
     *
     * @param {Point} point
     * @param {String|null} color
     * @returns {Point}
     */
    drawPoint(point, color) {
        if(color) {
            this.line.setColor(color);
            this._qDrawer.addToLine(point, color);
        } else {
            this._qDrawer.addToLine(point);
        }
        this.line.addPoint(point);

        return point;
    }

    /**
     * Canvas initialize
     */
    initialize() {
        this._canvas.width  = window.innerWidth - 20;
        this._canvas.height = window.innerHeight - 20;

        this._canvas.addEventListener('mousedown', (event) => {
            const point = new Point(event.pageX - this._canvas.offsetLeft, event.pageY - this._canvas.offsetTop);
            this.x = point.getX();
            this.y = point.getY();
            this.drawPoint(point);
            this._isDrawing = true;
        });

        this._canvas.addEventListener('mousemove', (event) => {
            if (this._isDrawing) {
                const point = new Point(event.pageX - this._canvas.offsetLeft, event.pageY - this._canvas.offsetTop);
                this.x = point.getX();
                this.y = point.getY();
                this.drawPoint(point);
            }
        }, false);

        this._canvas.addEventListener('mouseup', () => {
            this.x = 0;
            this.y = 0;
            this._isDrawing = false;
            this.connect.sendLine(this.line.serialize());
            this.line = new Line([], this.line.getColor());
            this._qDrawer.clearLine();
        });
    }
}