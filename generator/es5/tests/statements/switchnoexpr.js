$module('switchnoexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $t.fastbox(123, $g.____testlib.basictypes.Integer);
            if (false) {
              $current = 1;
              continue;
            } else {
              $current = 3;
              continue;
            }
            break;

          case 1:
            $t.fastbox(1234, $g.____testlib.basictypes.Integer);
            $current = 2;
            continue;

          case 2:
            $t.fastbox(789, $g.____testlib.basictypes.Integer);
            $resolve();
            return;

          case 3:
            if (true) {
              $current = 4;
              continue;
            } else {
              $current = 5;
              continue;
            }
            break;

          case 4:
            $t.fastbox(2345, $g.____testlib.basictypes.Integer);
            $current = 2;
            continue;

          case 5:
            if (true) {
              $current = 6;
              continue;
            } else {
              $current = 2;
              continue;
            }
            break;

          case 6:
            $t.fastbox(3456, $g.____testlib.basictypes.Integer);
            $current = 2;
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
