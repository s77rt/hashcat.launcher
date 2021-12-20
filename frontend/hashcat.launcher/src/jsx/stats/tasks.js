const tasks = {
	tasks: {},

	_update: function(taskUpdate) {
		var task = this.tasks[taskUpdate.task.id] || {
			id: taskUpdate.task.id,
			arguments: taskUpdate.task.arguments,
			process: taskUpdate.task.process,
			priority: taskUpdate.task.priority,
			stats: {},
			journal: []
		}

		// id is immutable
		// arguments are immutable
		task.process = taskUpdate.task.process
		task.priority = taskUpdate.task.priority

		try {
			task.stats = JSON.parse(taskUpdate.message)
		} catch (e) {
			if (taskUpdate.message.length > 0) {
				task.journal.push({
					message: taskUpdate.message,
					source: taskUpdate.source,
					timestamp: taskUpdate.timestamp
				})
			}
		}

		this.tasks[task.id] = task;
	},

	_delete: function(taskID) {
		delete(this.tasks[taskID]);
	}
}

export default tasks;
