import Point from "./Point.js";
import Line from "./Line.js";

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
        this._ownLine = new Line();
        window.requestAnimationFrame((timestamp) => {this.render(timestamp);});
    }

    /**
     * Add line to drawing query
     *
     * @param {Line} line
     */
    addQuery(line) {
        this._query.push(line);
    }

    /**
     * Add line to draw
     *
     * @param {Point} point
     * @param {String|null} color
     */
    addToLine(point, color) {
        this._ownLine.addPoint(point);
        if(color) {
            this._ownLine.setColor(color);
        }
    }

    /**
     * Number of current timestamp
     *
     * @param {Number} timestamp
     */
    render(timestamp) {
        window.requestAnimationFrame((timestamp) => {this.render(timestamp);});

        if(!this.start) {
            this.start = timestamp;
            return;
        }

        const distance = this._speed * (timestamp - this.start);
        this.start = timestamp;

        /**
         * @var {Line} line
         */
        this._query = this._query.map(line => {
            if(!line.isEnded()) {
                return this.drawDistance(line, distance);
            }
            return line;
        }).filter(line => {
            return !line.isEnded();
        });

        if(!this._ownLine.isEnded()) {
            this.drawDistance(this._ownLine, distance);
        }
    }

    /**
     * Clear line
     *
     * @param {String|null} color
     */
    clearLine(color) {
        if(!this._ownLine.isEnded()) {
            let line = this._ownLine;
            this._query.push(line);
        }
        this._ownLine = new Line([], color || '#000000');
    }

    /**
     * Drawing line per frame
     *
     * @param {Line} points
     * @param {Number} distance
     * @returns {Array}
     */
    drawDistance(line, distance) {
        const data = line.getDrawIteration();
        line.shift();
        const options = data.A.getAngleAndDistance(data.B);
        const newData = Point.calcNewPoint(data.A, options, distance);
        const newPoint = newData.point;

        this._context.beginPath();
        this._context.strokeStyle = 'black';
        this._context.lineWidth = 5;
        this._context.lineCap = "round";
        this._context.moveTo(newPoint.getX(), newPoint.getY());
        this._context.lineTo(data.A.getX(), data.A.getY());
        this._context.stroke();
        this._context.closePath();

        if(newData.distance > 0 && !line.isEnded()) {
            line = this.drawDistance(line, newData.distance);
        } else if(newData.distance < 0) {
            line.unshift(newPoint);
        }

        return line;
    }

    draw(from, to) {
        return [];
    }
}