const eventBus = {
	events: {},
	on(event, eventID, callback) {
		this.remove(event, eventID);

		if (!this.events[event])
			this.events[event] = {};

		this.events[event][eventID] = (e) => callback(e.detail);
		document.addEventListener(event, this.events[event][eventID]);
	},
	dispatch(event, data) {
		document.dispatchEvent(new CustomEvent(event, { detail: data }));
	},
	remove(event, eventID) {
		if (!this.events[event])
			return;

		if (this.events[event][eventID]) {
			document.removeEventListener(event, this.events[event][eventID]);
			delete(this.events[event][eventID]);
		}

		if (Object.keys(this.events[event]).length === 0)
			delete(this.events[event]);
	},
};

window.eventBus = eventBus;

export default eventBus;
