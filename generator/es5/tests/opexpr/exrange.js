$module('exrange', function () {
  var $static = this;
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
            $temp1 = $g.________testlib.basictypes.Integer.$exclusiverange($t.fastbox(0, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer));
            $current = 2;
            continue localasyncloop;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            i = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue localasyncloop;
            } else {
              $current = 5;
              continue localasyncloop;
            }
            break;

          case 4:
            counter = $t.fastbox(counter.$wrapped + i.$wrapped, $g.________testlib.basictypes.Integer);
            $current = 2;
            continue localasyncloop;

          case 5:
            $resolve($t.fastbox(counter.$wrapped == 1, $g.________testlib.basictypes.Boolean));
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
