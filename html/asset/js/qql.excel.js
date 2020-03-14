function getFileName(filePath) {
	var pattern = /\.{1}[a-z]{1,}$/;
    if (pattern.exec(filePath) !== null) {
        return (filePath.slice(0, pattern.exec(filePath).index));
    } else {
        return filePath;
    }
}

function uploadButtonEnabled(enabled) {
	$('#upload-file').toggleClass('disabled', !enabled);
}

$('#qql-csv-file').bind("change", function(e) {
	var files = e.target.files
	if (files.length > 0) {
		$('#download-card').attr('hidden', true)
		$('#xlsx-download').attr('href', "#");
		var file = files[0]
		if (file.type != 'text/csv') {
			alert('请上传csv格式文件')
			var inputFile = document.getElementById('qql-csv-file');
			inputFile.value = "";
			uploadButtonEnabled(false)
			return
		}

		if (file.size <= 0) {
			alert('文件内容为空，请选择其他文件')
			var inputFile = document.getElementById('qql-csv-file');
			inputFile.value = "";
			uploadButtonEnabled(false)
			return;
		}

		var fileName = file.name;
		var dateFormatterString = '2020-03-11-21-51-29';
		var targetFileName = getFileName(file.name);
		var date = targetFileName.substr(targetFileName.length - dateFormatterString.length, dateFormatterString.length);
		targetFileName = '柏兰德_' + date + '.xlsx';
		$('#xlsx-input-file').val(targetFileName);
		$('#xlsx-input-file').attr('placeholder', targetFileName)
		$('#xlsx-input-file').removeAttr('disabled')
		uploadButtonEnabled(true)
		$('#csv-file-name').html(fileName)
		$('#download-card').attr('hidden', true)
	}else {
		$('#csv-file-name').html('选择你的csv文件')
		$('#xlsx-input-file').attr('placeholder', '请输入文件名');
		$('#xlsx-input-file').val('')
		$('#xlsx-input-file').attr('disabled', 'disabled')
		uploadButtonEnabled(false)
	}
});

$('#xlsx-input-file').bind('input propertychange', function(e) {
	var elements = 'qql-csv-file';
	var fileName = $(this).val();
	var enabled = fileName.length >= 0;
	uploadButtonEnabled(enabled);
	if (enabled) {
		$('#xlsx-input-file').removeAttr('disabled')
	}else {
		$('#xlsx-input-file').attr('disabled', 'disabled')
	}
})


$("#upload-file").click(function(){
	var elements = 'qql-csv-file';
	var fileLists = document.getElementById(elements).files;
	if (fileLists.length <= 0 || typeof(fileLists) == 'undefined') {
		alert('csv文件为空')
		return;
	}

	var formData = new FormData();

	formData.append("csv_file",fileLists[0]);
	formData.append("fileName", $('#xlsx-input-file').val())

	$.ajax({
		url: '/qql/tool/excel/order',
		type: 'POST',
		processData: false,
		contentType: false,
		data: formData,
		success: function (s) {
			var obj = JSON.parse(s)
			$('#download-card').attr('hidden', !obj.success)
			if (!obj.success) {
				alert(obj.message)
				return;
			}
			var href = obj.data
			$('#xlsx-download').attr('href', href);
		}
	})
});

$("#download").click(function(){
	var link = document.getElementById('xlsx-download')
    link.click(); 
});

$("#clear-file").click(function(){
	$.ajax({
		url: '/qql/tool/excel/clearfiles',
		type: 'POST',
		success: function (s) {
			var obj = JSON.parse(s)
			$('#download-card').attr('hidden', true)
			if (obj.success) {
				alert('清除成功');
			}else {
				alert('文件清除失败:' + obj.message);
			}
		}
	})
});