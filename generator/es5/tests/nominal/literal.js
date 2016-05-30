$module('literal', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.box('foo' == 'foo', $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
});
