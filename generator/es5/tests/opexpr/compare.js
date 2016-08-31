$module('compare', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.$equals = function (first, second) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    };
    $static.$compare = function (first, second) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(1, $g.____testlib.basictypes.Integer));
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['compare', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Integer).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.compare.SomeClass).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $result;
    var first;
    var second;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.compare.SomeClass.new().then(function ($result0) {
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
            first = $result;
            $g.compare.SomeClass.new().then(function ($result0) {
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
            second = $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $t.box(!$t.unbox($result0), $g.____testlib.basictypes.Boolean);
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) < 0, $g.____testlib.basictypes.Boolean);
              $current = 5;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 5:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) > 0, $g.____testlib.basictypes.Boolean);
              $current = 6;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 6:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) <= 0, $g.____testlib.basictypes.Boolean);
              $current = 7;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 7:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) >= 0, $g.____testlib.basictypes.Boolean);
              $current = 8;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 8:
            $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $result0;
              $current = 9;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 9:
            $resolve($result);
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
