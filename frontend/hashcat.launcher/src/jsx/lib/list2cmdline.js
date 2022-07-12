export default function list2cmdline(seq) {
	var result = [];
	var needquote = false;
	seq.forEach((arg) => {
		var bs_buf = [];

        // Add a space to separate this argument from the others
		if (result.length > 0)
			result.push(' ');

		needquote = arg.includes(" ") || arg.includes("\t") || arg.length === 0;
		if (needquote)
			result.push('"');

		[...arg].forEach((c) => {
			if (c === '\\') {
				// Don't know if we need to double yet.
				bs_buf.push(c);
			} else if (c == '"') {
				// Double backslashes.
				result.push('\\'.repeat(bs_buf.length*2));
				bs_buf = [];
				result.push('\\"');
			} else {
				// Normal char
				if (bs_buf.length > 0) {
					result.push(...bs_buf);
					bs_buf = [];
				}
				result.push(c);
			}
		});

		if (bs_buf.length > 0)
			result.push(...bs_buf);

		if (needquote) {
			result.push(...bs_buf);
			result.push('"');
		}
	});

	return result.join('');
}
