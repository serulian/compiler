$module('matchexpr', function () {
  var $static = this;
  $static.DoSomething = function (someVar) {
    var $temp0;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.box(123, $g.____testlib.basictypes.Integer);
            $temp0 = someVar;
            $g.____testlib.basictypes.Integer.$equals($temp0, $t.box(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $t.unbox($result0);
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            if ($result) {
              $state.current = 2;
              continue;
            } else {
              $state.current = 4;
              continue;
            }
            break;

          case 2:
            $t.box(1234, $g.____testlib.basictypes.Integer);
            $state.current = 3;
            continue;

          case 3:
            $t.box(789, $g.____testlib.basictypes.Integer);
            $state.current = -1;
            return;

          case 4:
            $g.____testlib.basictypes.Integer.$equals($temp0, $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $t.unbox($result0);
              $state.current = 5;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 5:
            if ($result) {
              $state.current = 6;
              continue;
            } else {
              $state.current = 7;
              continue;
            }
            break;

          case 6:
            $t.box(2345, $g.____testlib.basictypes.Integer);
            $state.current = 3;
            continue;

          case 7:
            if (true) {
              $state.current = 8;
              continue;
            } else {
              $state.current = 3;
              continue;
            }
            break;

          case 8:
            $t.box(3456, $g.____testlib.basictypes.Integer);
            $state.current = 3;
            continue;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
