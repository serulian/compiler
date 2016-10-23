$module('varassign', function () {
  var $static = this;
  $static.DoSomething = function () {
    var i;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      i = $t.fastbox(2, $g.____testlib.basictypes.Integer);
      i = $t.fastbox(3, $g.____testlib.basictypes.Integer);
      $t.fastbox(1234, $g.____testlib.basictypes.Integer);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
