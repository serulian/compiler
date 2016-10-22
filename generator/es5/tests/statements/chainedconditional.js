$module('chainedconditional', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            if (false) {
              $current = 1;
              continue;
            } else {
              $current = 2;
              continue;
            }
            break;

          case 1:
            $t.fastbox(123, $g.____testlib.basictypes.Integer);
            $resolve($t.fastbox(false, $g.____testlib.basictypes.Boolean));
            return;

          case 2:
            if (false) {
              $current = 3;
              continue;
            } else {
              $current = 4;
              continue;
            }
            break;

          case 3:
            $t.fastbox(456, $g.____testlib.basictypes.Integer);
            $resolve($t.fastbox(false, $g.____testlib.basictypes.Boolean));
            return;

          case 4:
            $t.fastbox(789, $g.____testlib.basictypes.Integer);
            $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
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
