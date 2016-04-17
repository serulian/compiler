$module('loop', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $t.box(1234, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          case 1:
            $t.box(1357, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
