$module('lambda', function () {
  var $static = this;
  $static.DoSomething = $t.markpromising(function () {
    var $result;
    var $temp0;
    var $temp1;
    var l;
    var total;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            l = function () {
              var $current = 0;
              var $continue = function ($yield, $yieldin, $reject, $done) {
                while (true) {
                  switch ($current) {
                    case 0:
                      $yield($t.fastbox(1, $g.________testlib.basictypes.Integer));
                      $current = 1;
                      return;

                    case 1:
                      $yield($t.fastbox(2, $g.________testlib.basictypes.Integer));
                      $current = 2;
                      return;

                    default:
                      $done();
                      return;
                  }
                }
              };
              return $generator.new($continue, false, $t.any);
            };
            total = $t.fastbox(0, $g.________testlib.basictypes.Integer);
            $current = 1;
            continue localasyncloop;

          case 1:
            $promise.maybe(l()).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $temp1 = $result;
            $current = 3;
            continue localasyncloop;

          case 3:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            value = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 5;
              continue localasyncloop;
            } else {
              $current = 6;
              continue localasyncloop;
            }
            break;

          case 5:
            total = $t.fastbox(total.$wrapped + value.$wrapped, $g.________testlib.basictypes.Integer);
            $current = 3;
            continue localasyncloop;

          case 6:
            $resolve(total);
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
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.lambda.DoSomething()).then(function ($result0) {
              $result = $t.fastbox($result0.$wrapped == 3, $g.________testlib.basictypes.Boolean);
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
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
