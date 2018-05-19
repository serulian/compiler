$module('numeric', function () {
  var $static = this;
  $static.doSomething = $t.markpromising(function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.maybe((function () {
              return;
            })()).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve();
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $temp0;
    var $temp1;
    var counter;
    var i;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            counter = $t.fastbox(0, $g.________testlib.basictypes.Integer);
            $current = 1;
            continue localasyncloop;

          case 1:
            $temp1 = $g.________testlib.basictypes.Integer.$range($t.fastbox(0, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer));
            $current = 2;
            continue localasyncloop;

          case 2:
            $temp0 = $temp1.Next();
            i = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 3;
              continue localasyncloop;
            } else {
              $current = 5;
              continue localasyncloop;
            }
            break;

          case 3:
            counter = $t.fastbox(counter.$wrapped + i.$wrapped, $g.________testlib.basictypes.Integer);
            $promise.maybe($g.numeric.doSomething()).then(function ($result0) {
              $result = $result0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            $current = 2;
            continue localasyncloop;

          case 5:
            $resolve($t.fastbox(counter.$wrapped == 3, $g.________testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
