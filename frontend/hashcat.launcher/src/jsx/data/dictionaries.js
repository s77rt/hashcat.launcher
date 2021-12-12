import data from "./data"

export const getDictionaries = () => {
	return data.data.dictionaries || [];
};
