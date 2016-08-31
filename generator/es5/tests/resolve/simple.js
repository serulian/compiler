$module('simple', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      a = $t.box(true, $g.____testlib.basictypes.Boolean);
      $resolve(a);
      return;
    };
    return $promise.new($continue);
  };
});
