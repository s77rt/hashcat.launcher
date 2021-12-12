const eventBus = {
	events: {},
	on(event, callback) {
		this.events[event] = (e) => callback(e.detail);
		document.addEventListener(event, this.events[event]);
	},
	dispatch(event, data) {
		document.dispatchEvent(new CustomEvent(event, { detail: data }));
	},
	remove(event) {
		document.removeEventListener(event, this.events[event]);
		delete(this.events[event]);
	},
};

window.eventBus = eventBus;

export default eventBus;
