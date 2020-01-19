export default class Serializer {
    /**
     * Serialize object to json
     *
     * @param {Object} obj
     */
    static serializeObjectToJson(obj) {
        return obj.serialize();
    }

    /**
     * Serialize array to json
     *
     * @param {Array} array
     */
    static serializeArrayToJson(array) {
        let json = "[";
        array.map(obj => {
            json += obj.serialize();
        });
        json += "]";

        return json;
    }
}