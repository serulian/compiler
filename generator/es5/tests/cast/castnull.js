$module('castnull', function () {
  var $static = this;
  $static.TEST = function () {
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      value = null;
      $resolve($t.cast(value, $g.____testlib.basictypes.Boolean, false));
      return;
    };
    return $promise.new($continue);
  };
});
