$module('varnoinit', function () {
  var $static = this;
  $static.DoSomething = function () {
    var i;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.box(1234, $g.____testlib.basictypes.Integer);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
