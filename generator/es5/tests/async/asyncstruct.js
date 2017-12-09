$module('asyncstruct', function () {
  var $static = this;
  this.$struct('4e1a7940', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Foo, Bar) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Foo: Foo,
        Bar: Bar,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Foo', 'Foo', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'Bar', 'Bar', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|0b2e6e78<4e1a7940>": true,
        "equals|4|0b2e6e78<5e61c39d>": true,
        "Stringify|2|0b2e6e78<c509e19d>": true,
        "Mapping|2|0b2e6e78<204295f9<any>>": true,
        "Clone|2|0b2e6e78<4e1a7940>": true,
        "String|2|0b2e6e78<c509e19d>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('dae69f8e', function (s) {
    return $g.________testlib.basictypes.String.$plus(s.Foo.String(), s.Bar.String());
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var vle;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.asyncstruct.DoSomethingAsync($g.asyncstruct.SomeStruct.new($t.fastbox(1, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer)))).then(function ($result0) {
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
            vle = $result;
            $resolve($g.________testlib.basictypes.String.$equals(vle, $t.fastbox("12", $g.________testlib.basictypes.String)));
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
