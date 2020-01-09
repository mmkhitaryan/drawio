import Point from "./Point";

export default class QueryDrawer {
    /**
     * QueryDreawer constructor
     * @param {Node} canvas
     * @param {Number} speed
     */
    constructor(canvas, speed) {
        this.start = 0;
        this._speed = speed;
        this._context = canvas.getContext("2d");
        this._query = [];
        this._ownLine = [];
        window.requestAnimationFrame(this.render);
    }

    /**
     * Add line to drawing query
     *
     * @param {Array} points
     */
    addQuery(points) {
        this._query.push(points);
    }

    /**
     * Add line to drawing query
     *
     * @param {Uint16Array} points
     */
    addToLine(point) {
        this._ownLine.push(point);
    }

    /**
     * Number of current timestamp
     *
     * @param {Number} timestamp
     */
    render(timestamp) {
        window.requestAnimationFrame(this.render);

        if(!this.start) {
            this.start = timestamp;
            return;
        }

        const distance = speed * (timestamp - this.start);
        this.start = timestamp;

        this._query.map(points => {
            return this.drawDistance(points, distance);
        }).filter(points => {
            return points.length > 1;
        });

        if(this._ownLine.length > 1) {
            this.drawDistance(this._ownLine);
        }
    }

    clearLine() {
        this._ownLine = [];
    }

    /**
     * Drawing line per frame
     *
     * @param {Array} points
     * @param {Number} distance
     * @returns {Array}
     */
    drawDistance(points, distance) {
        const oldPoint = Point.decodePoint(points[0]);
        let newPoint = Point.decodePoint(points[1]);
        const options = Point.getAngleAndDistance(oldPoint, newPoint);
        const data = Point.calcNewPoint(oldPoint, options, distance);
        newPoint = data.point;

        this._context.beginPath();
        this._context.strokeStyle = 'black';
        this._context.lineWidth = 5;
        this._context.lineCap = "round";
        this._context.moveTo(newPoint.x, newPoint.y);
        this._context.lineTo(oldPoint.x, oldPoint.y);
        this._context.stroke();
        this._context.closePath();
        points.shift();

        if(data.distance > 0 && points.length > 1) {
            points = this.drawDistance(points, data.distance);
        } else if(data.distance < 0) {
            points.unshift(data.newPoint);
        }

        return points;
    }
}