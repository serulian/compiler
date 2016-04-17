$module('varassign', function () {
  var $static = this;
  $static.DoSomething = function () {
    var i;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      i = $t.box(2, $g.____testlib.basictypes.Integer);
      i = $t.box(3, $g.____testlib.basictypes.Integer);
      $t.box(1234, $g.____testlib.basictypes.Integer);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
