$module('loopexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $t.fastbox(1234, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          case 1:
            if (true) {
              $current = 2;
              continue;
            } else {
              $current = 3;
              continue;
            }
            break;

          case 2:
            $t.fastbox(1357, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          case 3:
            $t.fastbox(5678, $g.____testlib.basictypes.Integer);
            $resolve();
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
