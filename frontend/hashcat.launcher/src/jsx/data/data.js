const data = {
	data: {},
	callback: undefined,

	getHashes: async function() {
		if (typeof window.GOgetHashes !== "function") {
			console.error("GOgetHashes is not a function");
			return
		}
		window.GOgetHashes().then(
			response => {
				this.data.hashes = response;
				if (typeof this.callback === "function") {
					this.callback();
				}
			},
			error => {
				console.error("Failed to get hashes", error);
			}
		);
	},

	getAlgorithms: async function() {
		if (typeof window.GOgetAlgorithms !== "function") {
			console.error("GOgetAlgorithms is not a function");
			return
		}
		window.GOgetAlgorithms().then(
			response => {
				this.data.algorithms = response;
				if (typeof this.callback === "function") {
					this.callback();
				}
			},
			error => {
				console.error("Failed to get algorithms", error);
			}
		);
	},

	getDictionaries: async function() {
		if (typeof window.GOgetDictionaries !== "function") {
			console.error("GOgetDictionaries is not a function");
			return
		}
		window.GOgetDictionaries().then(
			response => {
				this.data.dictionaries = response;
				if (typeof this.callback === "function") {
					this.callback();
				}
			},
			error => {
				console.error("Failed to get dictionaries", error);
			}
		);
	},

	getRules: async function() {
		if (typeof window.GOgetRules !== "function") {
			console.error("GOgetRules is not a function");
			return
		}
		window.GOgetRules().then(
			response => {
				this.data.rules = response;
				if (typeof this.callback === "function") {
					this.callback();
				}
			},
			error => {
				console.error("Failed to get rules", error);
			}
		);
	},

	getMasks: async function() {
		if (typeof window.GOgetMasks !== "function") {
			console.error("GOgetMasks is not a function");
			return
		}
		window.GOgetMasks().then(
			response => {
				this.data.masks = response;
				if (typeof this.callback === "function") {
					this.callback();
				}
			},
			error => {
				console.error("Failed to get masks", error);
			}
		);
	}
}

window.data = data;

export default data;
