
$('#todoCheck1').TodoList({
	onCheck: function(checkbox) {
		alert('勾选了');
	},
	onUnCheck: function(checkbox) {
		alert('未勾选');
	}
});


$('#my-todo-list').TodoList({
  onCheck: function(checkbox) {
    // Do something when the checkbox is checked
  },
  onUnCheck: function(checkbox) {
    // Do something after the checkbox has been unchecked
  }
})


$(function(){
	$('#content-frame').load('/html/qql.html');
})