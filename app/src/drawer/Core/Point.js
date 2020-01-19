export default class Point {
    /**
     * Point constructor
     * @param {Number} x
     * @param {Number} y
     */
    constructor(x, y) {
        this._x = x;
        this._y = y;
    }

    /**
     * Get point X
     */
    getX() {
        return this._x;
    }

    /**
     * Get point Y
     */
    getY() {
        return this._y;
    }

    /**
     * Get angle and distance by start and end points
     *
     * @param {Point} point
     * @returns {Object}
     */
     getAngleAndDistance(point) {
        return {
            angle:  Math.atan2(point.getY() - this._y, point.getX() - this._x) * (180/Math.PI),
            distance: Math.sqrt(Math.pow((point.getY() - this._y), 2) + Math.pow((point.getX() - this._x), 2))
        }
    }

    /**
     * Get endpoint of line by angle and distance
     *
     * @param {Number} angle
     * @param {Number} distance
     * @returns {Object}
     */
    getEndPoint(angle, distance) {
        const radians = angle/(180/Math.PI);

        return new Point((Math.cos(radians) * distance) + this._x, (Math.sin(radians) * distance) + this._y);
    }

    /**
     * Encode point to uint16 array
     *
     * @returns {Uint16Array}
     */
    getEncodedPoint() {
        return new Uint16Array([
            parseInt(this._x, 10),
            parseInt(this._y, 10)
        ]);
    }

    /**
     * Calc distance to move brush
     *
     * @param {Point} point
     * @param {Object} options
     * @param {Number} distance
     * @return {Object}
     */
    static calcNewPoint(point, options, distance)
    {
        const least = distance - options.distance;
        point = point.getEndPoint(options.angle, least < 0 ? distance : options.distance);

        return {
            point: point,
            distance: least
        };
    }
}