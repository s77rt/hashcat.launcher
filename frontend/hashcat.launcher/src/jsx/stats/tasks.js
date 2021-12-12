const tasks = {
	tasks: {},

	_update: function(taskUpdate) {
		var task = this.tasks[taskUpdate.task.id] || {
			id: taskUpdate.task.id,
			arguments: taskUpdate.task.arguments,
			priority: taskUpdate.task.priority,
			stats: {},
			journal: []
		}

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
}

export default tasks;
