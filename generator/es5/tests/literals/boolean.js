$module('boolean', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.box(true, $g.____testlib.basictypes.Boolean);
      $t.box(false, $g.____testlib.basictypes.Boolean);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
