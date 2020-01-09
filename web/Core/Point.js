export default class Point {
    /**
     * Get angle and distance by start and end points
     *
     * @param {Object} a
     * @param {Object} b
     * @returns {Object}
     */
    static getAngleAndDistance(a, b) {
        return {
            angle:  Math.atan2(b.y - a.y, b.x - a.x) * (180/Math.PI),
            distance: Math.sqrt(Math.pow((b.y - a.y), 2) + Math.pow((b.x - a.x), 2))
        }
    }

    /**
     * Get endpoint of line by angle and distance
     *
     * @param {Object} a
     * @param {Number} angle
     * @param {Number} distance
     * @returns {Object}
     */
    static getEndPoint(a, angle, distance) {
        const radians = angle/(180/Math.PI);
        return {
            x: (Math.cos(radians) * distance) + a.x,
            y: (Math.sin(radians) * distance) + a.y
        }
    }

    /**
     * Encode point to uint16 array
     *
     * @param {Object} point
     * @returns {Uint16Array}
     */
    static encodePoint(point) {
        return new Uint16Array([
            parseInt(point.x, 10),
            parseInt(point.y, 10)
        ]);
    }

    /**
     * Decode uint16 array to object
     *
     * @param {Uint16Array} encoded
     * @returns {Object}
     */
    static decodePoint(encoded) {
        return {
            x: encoded[0],
            y: encoded[1]
        };
    }

    /**
     * Calc distance to move brush
     *
     * @param {Object} point
     * @param Object} options
     * @param {Number} distance
     * @return {Object}
     */
    static calcNewPoint(point, options, distance)
    {
        //TODO: Write moving algorithm
        return {
            point: point,
            distance: distance
        };
    }
}