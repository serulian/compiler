$module('inlineloop', function () {
  var $static = this;
  $static.somenumber = function (props, value) {
    return value;
  };
  $static.somesum = $t.markpromising(function (props, numbers) {
    var $result;
    var $temp0;
    var $temp1;
    var number;
    var sum;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            sum = $t.fastbox(0, $g.________testlib.basictypes.Integer);
            $current = 1;
            continue localasyncloop;

          case 1:
            $temp1 = numbers;
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
            number = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue localasyncloop;
            } else {
              $current = 5;
              continue localasyncloop;
            }
            break;

          case 4:
            sum = $t.fastbox(sum.$wrapped + number.$wrapped, $g.________testlib.basictypes.Integer);
            $current = 2;
            continue localasyncloop;

          case 5:
            $resolve(sum);
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
    var result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.inlineloop.somesum($g.________testlib.basictypes.Mapping($t.any).Empty(), $g.________testlib.basictypes.MapStream($g.________testlib.basictypes.Integer, $g.________testlib.basictypes.Integer)($g.________testlib.basictypes.Integer.$range($t.fastbox(0, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer)), function (index) {
              return $g.inlineloop.somenumber($g.________testlib.basictypes.Mapping($t.any).Empty(), index);
            }))).then(function ($result0) {
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
            result = $result;
            $resolve($t.fastbox(result.$wrapped == 3, $g.________testlib.basictypes.Boolean));
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
