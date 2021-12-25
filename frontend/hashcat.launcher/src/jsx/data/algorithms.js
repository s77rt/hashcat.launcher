import data from "./data"

export const getAlgorithms = () => {
	return data.data.algorithms || {};
};
