$module('switchexpr', function () {
  var $static = this;
  $static.DoSomething = function (someVar) {
    var $result;
    var $temp0;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $t.fastbox(123, $g.____testlib.basictypes.Integer);
            $temp0 = someVar;
            $g.____testlib.basictypes.Integer.$equals($temp0, $t.fastbox(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0.$wrapped;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            if ($result) {
              $current = 2;
              continue;
            } else {
              $current = 4;
              continue;
            }
            break;

          case 2:
            $t.fastbox(1234, $g.____testlib.basictypes.Integer);
            $current = 3;
            continue;

          case 3:
            $t.fastbox(789, $g.____testlib.basictypes.Integer);
            $resolve();
            return;

          case 4:
            $g.____testlib.basictypes.Integer.$equals($temp0, $t.fastbox(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0.$wrapped;
              $current = 5;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 5:
            if ($result) {
              $current = 6;
              continue;
            } else {
              $current = 7;
              continue;
            }
            break;

          case 6:
            $t.fastbox(2345, $g.____testlib.basictypes.Integer);
            $current = 3;
            continue;

          case 7:
            if (true) {
              $current = 8;
              continue;
            } else {
              $current = 3;
              continue;
            }
            break;

          case 8:
            $t.fastbox(3456, $g.____testlib.basictypes.Integer);
            $current = 3;
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
