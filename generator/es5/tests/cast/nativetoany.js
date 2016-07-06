$module('nativetoany', function () {
  var $static = this;
  $static.TEST = function () {
    var b;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      b = true;
      $resolve($t.cast($t.cast(b, $t.any), $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
});
