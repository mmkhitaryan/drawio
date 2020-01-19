import Point from './Point.js';

export default class Line {
    constructor(points, color) {
        this._points = points ? points : [];
        this._color = color || '#000000';
    }

    /**
     * Get line color
     */
    getColor() {
        return this._color;
    }

    /**
     * Set line color
     * @param {String|null} color
     */
    setColor(color) {
        this._color = color || '#000000';

        return this;
    }

    /**
     * Add new point to line
     * @param {Point} point
     */
    addPoint(point) {
        this._points.push(point);

        return this;
    }

    /**
     * Shift points array
     */
    shift() {
        this._points.shift();

        return this;
    }

    /**
     * Unshift points array
     * @param {Point} point
     */
    unshift(point) {
        this._points.unshift(point);

        return this;
    }

    length() {
        return this._points.length;
    }

    /**
     * Clear all line points
     */
    clearLine() {
        this._points = [];

        return this;
    }

    /**
     * Get next draw iteration points&color
     * @returns {null|Object}
     */
    getDrawIteration() {
        if(!this.isEnded()) {
            return {
                A: this._points[0],
                B: this._points[1],
                color: this._color
            };
        }
        return null;
    }

    /**
     * If line is ended
     */
    isEnded() {
        return this._points.length < 2;
    }

    /**
     * Serialize object to json
     */
    serialize() {
        let data = {
            points: [],
            color: this._color
        };
        this._points.map(point => {
            data.points.push([point.getX(), point.getY()]);
        });
        return JSON.stringify(data);
    }
}